package lnd

import (
	"context"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
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
	lndMinPoolCapacity      = 6
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
			return nil
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
			return nil
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
	clt, err := l.Client(ctx)
	if err != nil {
		return err
	}
	defer clt.Close()

	channelEvents, err := clt.SubscribeChannelEvents(ctx, &lnrpc.ChannelEventSubscription{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			event, err := channelEvents.Recv()
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.Canceled {
					l.logger.Debug("stopping subscribe channels: context canceled")
					return nil
				}
				return err
			}
			if event.Type == lnrpc.ChannelEventUpdate_FULLY_RESOLVED_CHANNEL {
				events <- &models.ChannelUpdate{}
			}

		}
	}
}

func chanpointToString(c *lnrpc.ChannelPoint) string {
	hash := c.GetFundingTxidBytes()
	for i := 0; i < len(hash)/2; i++ {
		hash[i], hash[len(hash)-i-1] = hash[len(hash)-i-1], hash[i]
	}
	output := c.OutputIndex
	result := fmt.Sprintf("%s:%d", hex.EncodeToString(hash), output)
	return result
}

func (l Backend) SubscribeGraphEvents(ctx context.Context, events chan *models.ChannelEdgeUpdate) error {
	clt, err := l.Client(ctx)
	if err != nil {
		return err
	}
	defer clt.Close()

	graphEvents, err := clt.SubscribeChannelGraph(ctx, &lnrpc.GraphTopologySubscription{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			event, err := graphEvents.Recv()
			if err != nil {
				st, ok := status.FromError(err)
				if ok && st.Code() == codes.Canceled {
					l.logger.Debug("stopping subscribe graph: context canceled")
					return nil
				}
				return err
			}
			chanPoints := []string{}
			for _, c := range event.ChannelUpdates {
				chanPoints = append(chanPoints, chanpointToString(c.ChanPoint))
			}
			if len(chanPoints) > 0 {
				events <- &models.ChannelEdgeUpdate{ChanPoints: chanPoints}
			}
		}
	}
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
			return nil
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
	channel.LocalPolicy = protoToRoutingPolicy(resp.Node1Policy)
	channel.RemotePolicy = protoToRoutingPolicy(resp.Node2Policy)

	info, err := clt.GetInfo(ctx, &lnrpc.GetInfoRequest{})
	if err != nil {
		return errors.WithStack(err)
	}
	if info != nil && resp.Node1Pub != info.IdentityPubkey {
		channel.LocalPolicy, channel.RemotePolicy = channel.RemotePolicy, channel.LocalPolicy
	}

	return nil
}

func (l Backend) GetNode(ctx context.Context, pubkey string, includeChannels bool) (*models.Node, error) {
	l.logger.Debug("GetNode")

	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()

	req := &lnrpc.NodeInfoRequest{PubKey: pubkey, IncludeChannels: includeChannels}
	resp, err := clt.GetNodeInfo(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result := nodeProtoToNode(resp)
	if forcedAlias, ok := l.cfg.Aliases[result.PubKey]; ok {
		result.ForcedAlias = forcedAlias
	}
	return result, nil
}

func (l Backend) GetForwardingHistory(ctx context.Context, startTime string, maxNumEvents uint32) ([]*models.ForwardingEvent, error) {
	l.logger.Debug("GetForwardingHistory")

	clt, err := l.Client(ctx)
	if err != nil {
		return nil, err
	}
	defer clt.Close()
	t, err := parseTime(startTime, time.Now())
	req := &lnrpc.ForwardingHistoryRequest{
		StartTime:    t,
		NumMaxEvents: maxNumEvents,
	}
	resp, err := clt.ForwardingHistory(ctx, req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	result := protoToForwardingHistory(resp)

	// Enrich peer alias names.
	// This can be removed once the ForwardingHistory
	// contains the peer aliases by default.
	enrichPeerAliases := func(ctx context.Context, events []*models.ForwardingEvent) error {

		if len(events) == 0 {
			return nil
		}

		selfInfo, err := clt.GetInfo(ctx, &lnrpc.GetInfoRequest{})
		if err != nil {
			return errors.WithStack(err)
		}

		getPeerAlias := func(chanId uint64) (string, error) {
			chanInfo, err := clt.GetChanInfo(ctx, &lnrpc.ChanInfoRequest{
				ChanId: chanId,
			})
			if err != nil {
				return "", errors.WithStack(err)
			}
			pubKey := chanInfo.Node1Pub
			if selfInfo.IdentityPubkey == chanInfo.Node1Pub {
				pubKey = chanInfo.Node2Pub
			}
			nodeInfo, err := clt.GetNodeInfo(ctx, &lnrpc.NodeInfoRequest{
				PubKey: pubKey,
			})
			if err != nil {
				return "", errors.WithStack(err)
			}

			return nodeInfo.Node.Alias, nil
		}

		cache := make(map[uint64]string)
		for i, event := range events {

			if val, ok := cache[event.ChanIdIn]; ok {
				events[i].PeerAliasIn = val
			} else {
				events[i].PeerAliasIn, err = getPeerAlias(event.ChanIdIn)
				if err != nil {
					cache[event.ChanIdIn] = events[i].PeerAliasIn
				}
			}

			if val, ok := cache[event.ChanIdOut]; ok {
				events[i].PeerAliasOut = val
			} else {
				events[i].PeerAliasOut, err = getPeerAlias(event.ChanIdOut)
				if err != nil {
					cache[event.ChanIdOut] = events[i].PeerAliasOut
				}
			}
		}

		return nil

	}
	err = enrichPeerAliases(ctx, result)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
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

// reTimeRange matches systemd.time-like short negative timeranges, e.g. "-200s".
var reTimeRange = regexp.MustCompile(`^-\d{1,18}[s|m|h|d|w|M|y]$`)

// secondsPer allows translating s(seconds), m(minutes), h(ours), d(ays),
// w(eeks), M(onths) and y(ears) into corresponding seconds.
var secondsPer = map[string]int64{
	"s": 1,
	"m": 60,
	"h": 3600,
	"d": 86400,
	"w": 604800,
	"M": 2630016,  // 30.44 days
	"y": 31557600, // 365.25 days
}

// parseTime parses UNIX timestamps or short timeranges inspired by systemd
// (when starting with "-"), e.g. "-1M" for one month (30.44 days) ago.
func parseTime(s string, base time.Time) (uint64, error) {
	if reTimeRange.MatchString(s) {
		last := len(s) - 1

		d, err := strconv.ParseInt(s[1:last], 10, 64)
		if err != nil {
			return uint64(0), err
		}

		mul := secondsPer[string(s[last])]
		return uint64(base.Unix() - d*mul), nil
	}

	return strconv.ParseUint(s, 10, 64)
}
