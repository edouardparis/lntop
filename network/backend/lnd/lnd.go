package lnd

import (
	"context"
	"fmt"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network/backend/pool"
	"github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/network/options"
)

const (
	lndDefaultInvoiceExpiry = 3600
	lndMinPoolCapacity      = 4
)

type Client struct {
	lnrpc.LightningClient
	conn *pool.Conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}

type RouterClient struct {
	routerrpc.RouterClient
	conn *pool.Conn
}

func (c *RouterClient) Close() error {
	return c.conn.Close()
}

type Backend struct {
	cfg    *config.Network
	logger logging.Logger
	pool   *pool.Pool
}

func (l Backend) NodeName() string {
	return l.cfg.Name
}

func (l Backend) Ping() error {
	clt, err := l.Client(context.Background())
	if err != nil {
		return err
	}
	defer clt.Close()
	return nil
}

func (l Backend) Info(ctx context.Context) (*models.Info, error) {
	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	resp, err := clt.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		return nil, err
	}

	return infoProtoToInfo(resp), nil
}

func (l Backend) SubscribeInvoice(ctx context.Context, channelInvoice chan *models.Invoice) error {
	clt, err := l.Client(ctx)
	if err != nil {
		return err
	}
	defer clt.Close()

	cltInvoices, err := clt.SubscribeInvoices(ctx, &lnrpc.InvoiceSubscription{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			break
		default:
			invoice, err := cltInvoices.Recv()
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.Canceled {
					l.logger.Debug("stopping subscribe invoice: context canceled")
					return nil
				}
				return err
			}

			channelInvoice <- lookupInvoiceProtoToInvoice(invoice)
		}
	}
}

func (l Backend) SubscribeTransactions(ctx context.Context, channel chan *models.Transaction) error {
	clt, err := l.Client(ctx)
	if err != nil {
		return err
	}
	defer clt.Close()

	cltTransactions, err := clt.SubscribeTransactions(ctx, &lnrpc.GetTransactionsRequest{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			break
		default:
			transaction, err := cltTransactions.Recv()
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.Canceled {
					l.logger.Debug("stopping subscribe transactions: context canceled")
					return nil
				}
				return err
			}

			channel <- protoToTransaction(transaction)
		}
	}
}

func (l Backend) SubscribeChannels(ctx context.Context, events chan *models.ChannelUpdate) error {
	_, err := l.Client(ctx)
	if err != nil {
		return err
	}

	// events, err := clt.SubscribeChannelEvents(ctx, &lnrpc.ChannelEventSubscription{})
	// if err != nil {
	// 	return err
	// }

	// for {
	// 	event, err := events.Recv()
	// 	if err != nil {
	// 		return err
	// 	}
	// events <-
	//}
	return nil
}

func (l Backend) SubscribeRoutingEvents(ctx context.Context, channelEvents chan *models.RoutingEvent) error {
	clt, err := l.RouterClient(ctx)
	if err != nil {
		return err
	}
	defer clt.Close()

	cltRoutingEvents, err := clt.SubscribeHtlcEvents(ctx, &routerrpc.SubscribeHtlcEventsRequest{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			break
		default:
			event, err := cltRoutingEvents.Recv()
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.Canceled {
					l.logger.Debug("stopping subscribe routing events: context canceled")
					return nil
				}
				return err
			}

			channelEvents <- protoToRoutingEvent(event)
		}
	}
}

func (l Backend) Client(ctx context.Context) (*Client, error) {
	conn, err := l.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		LightningClient: lnrpc.NewLightningClient(conn.ClientConn),
		conn:            conn,
	}, nil
}

func (l Backend) RouterClient(ctx context.Context) (*RouterClient, error) {
	conn, err := l.pool.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &RouterClient{
		RouterClient: routerrpc.NewRouterClient(conn.ClientConn),
		conn:         conn,
	}, nil
}

func (l Backend) NewClientConn() (*grpc.ClientConn, error) {
	return newClientConn(l.cfg)
}

func (l Backend) GetTransactions(ctx context.Context) ([]*models.Transaction, error) {
	l.logger.Debug("Get transactions...")
	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	req := &lnrpc.GetTransactionsRequest{}
	resp, err := clt.GetTransactions(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return protoToTransactions(resp), nil
}

func (l Backend) GetWalletBalance(ctx context.Context) (*models.WalletBalance, error) {
	l.logger.Debug("Retrieve wallet balance...")

	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	req := &lnrpc.WalletBalanceRequest{}
	resp, err := clt.WalletBalance(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	balance := protoToWalletBalance(resp)

	l.logger.Debug("Wallet balance retrieved", logging.Object("wallet", balance))

	return balance, nil
}

func (l Backend) GetChannelsBalance(ctx context.Context) (*models.ChannelsBalance, error) {
	l.logger.Debug("Retrieve channel balance...")

	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	req := &lnrpc.ChannelBalanceRequest{}
	resp, err := clt.ChannelBalance(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	balance := protoToChannelsBalance(resp)

	l.logger.Debug("Channel balance retrieved", logging.Object("balance", balance))

	return balance, nil
}

func (l Backend) ListChannels(ctx context.Context, opt ...options.Channel) ([]*models.Channel, error) {
	l.logger.Debug("List channels")

	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	opts := options.NewChannelOptions(opt...)
	req := &lnrpc.ListChannelsRequest{
		ActiveOnly:   opts.Active,
		InactiveOnly: opts.Inactive,
		PublicOnly:   opts.Public,
		PrivateOnly:  opts.Private,
	}

	resp, err := clt.ListChannels(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	channels := listChannelsProtoToChannels(resp)

	if opts.Pending {
		req := &lnrpc.PendingChannelsRequest{}
		resp, err := clt.PendingChannels(ctx, req)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		channels = append(channels, pendingChannelsProtoToChannels(resp)...)
	}

	fields := make([]logging.Field, len(channels))
	for i := range channels {
		fields[i] = logging.Object(fmt.Sprintf("channel_%d", i), channels[i])
	}

	l.logger.Debug("Channels retrieved", fields...)

	return channels, nil
}

func (l Backend) GetChannelInfo(ctx context.Context, channel *models.Channel) error {
	l.logger.Debug("GetChannelInfo")

	// If channel does not have ID (pending), information cannot be retrieved
	if channel.ID == 0 {
		return nil
	}

	clt, err := l.Client(ctx)
	if err != nil {
		return err
	}
	defer clt.Close()

	req := &lnrpc.ChanInfoRequest{ChanId: channel.ID}
	resp, err := clt.GetChanInfo(ctx, req)
	if err != nil {
		return errors.WithStack(err)
	}
	if resp == nil {
		return nil
	}

	t := time.Unix(int64(uint64(resp.LastUpdate)), 0)
	channel.LastUpdate = &t
	channel.Policy1 = protoToRoutingPolicy(resp.Node1Policy)
	channel.Policy2 = protoToRoutingPolicy(resp.Node2Policy)

	return nil
}

func (l Backend) GetNode(ctx context.Context, pubkey string) (*models.Node, error) {
	l.logger.Debug("GetNode")

	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	req := &lnrpc.NodeInfoRequest{PubKey: pubkey}
	resp, err := clt.GetNodeInfo(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return nodeProtoToNode(resp), nil
}

func (l Backend) CreateInvoice(ctx context.Context, amount int64, desc string) (*models.Invoice, error) {
	l.logger.Debug("Create invoice...",
		logging.Int64("amount", amount),
		logging.String("desc", desc))

	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	req := &lnrpc.Invoice{
		Value:        amount,
		Memo:         desc,
		CreationDate: time.Now().Unix(),
		Expiry:       lndDefaultInvoiceExpiry,
	}

	resp, err := clt.AddInvoice(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	invoice := addInvoiceProtoToInvoice(req, resp)

	l.logger.Debug("Invoice retrieved", logging.Object("invoice", invoice))

	return invoice, nil
}

func (l Backend) GetInvoice(ctx context.Context, RHash string) (*models.Invoice, error) {
	l.logger.Debug("Retrieve invoice...", logging.String("r_hash", RHash))

	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	req := &lnrpc.PaymentHash{
		RHashStr: RHash,
	}

	resp, err := clt.LookupInvoice(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	invoice := lookupInvoiceProtoToInvoice(resp)

	l.logger.Debug("Invoice retrieved", logging.Object("invoice", invoice))

	return invoice, nil
}

func (l Backend) SendPayment(ctx context.Context, payreq *models.PayReq) (*models.Payment, error) {
	l.logger.Debug("Send payment...",
		logging.String("destination", payreq.Destination),
		logging.Int64("amount", payreq.Amount),
	)

	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	req := &lnrpc.SendRequest{PaymentRequest: payreq.String}

	resp, err := clt.SendPaymentSync(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	payment := sendPaymentProtoToPayment(payreq, resp)

	l.logger.Debug("Payment paid", logging.Object("payment", payment))

	return payment, nil
}

func (l Backend) DecodePayReq(ctx context.Context, payreq string) (*models.PayReq, error) {
	l.logger.Info("decode payreq", logging.String("payreq", payreq))
	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	resp, err := clt.DecodePayReq(ctx, &lnrpc.PayReqString{PayReq: payreq})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return payreqProtoToPayReq(resp, payreq), nil
}

func New(c *config.Network, logger logging.Logger) (*Backend, error) {
	var err error

	backend := &Backend{
		cfg:    c,
		logger: logger.With(logging.String("name", c.Name)),
	}

	if c.PoolCapacity < lndMinPoolCapacity {
		c.PoolCapacity = lndMinPoolCapacity
		logger.Info("pool_capacity too small, ignoring")
	}
	backend.pool, err = pool.New(backend.NewClientConn, c.PoolCapacity, time.Duration(c.ConnTimeout))
	if err != nil {
		return nil, err
	}

	return backend, nil
}
