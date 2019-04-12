package lightningd

import (
	"context"
	"time"

	"github.com/fiatjaf/lightningd-gjson-rpc"
	"github.com/tidwall/gjson"

	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/network/options"
)

type Backend struct {
	client *lightning.Client
	cfg    *config.Network
	logger logging.Logger
}

func New(c *config.Network, logger logging.Logger) (*Backend, error) {
	backend := &Backend{
		cfg:    c,
		logger: logger.With(logging.String("name", c.Name)),
	}

	backend.client = &lightning.Client{Path: c.SocketPath}
	return backend, nil
}

func (l Backend) NodeName() string {
	return l.cfg.Name
}

func (l Backend) Ping() error {
	_, err := l.client.Call("help")
	return err
}

func (l Backend) Info(ctx context.Context) (*models.Info, error) {
	res, err := l.client.Call("getinfo")
	if err != nil {
		return nil, err
	}

	return &models.Info{
		res.Get("id").String(),
		res.Get("alias").String(),
		uint32(res.Get("num_pending_channels").Int()),
		uint32(res.Get("num_active_channels").Int()),
		uint32(res.Get("num_inactive_channels").Int()),
		uint32(res.Get("num_peers").Int()),
		uint32(res.Get("blockheight").Int()),
		"",
		true,
		res.Get("version").String(),
		[]string{"bitcoin"},
		res.Get("network").String() != "bitcoin",
	}, nil
}

func (l Backend) SubscribeInvoice(ctx context.Context, channelInvoice chan *models.Invoice) error {
	return nil
}

func (l Backend) SubscribeChannels(ctx context.Context, events chan *models.ChannelUpdate) error {
	return nil
}

func (l Backend) GetWalletBalance(ctx context.Context) (*models.WalletBalance, error) {
	resp, err := l.client.Call("listfunds")
	if err != nil {
		return nil, err
	}

	wb := models.WalletBalance{}

	resp.Get("outputs").ForEach(func(_, utxo gjson.Result) bool {
		amount := utxo.Get("value").Int()

		if utxo.Get("status").String() == "confirmed" {
			wb.ConfirmedBalance += amount
		} else {
			wb.UnconfirmedBalance += amount
		}
		wb.TotalBalance += amount
		return true
	})

	return &wb, nil
}

func (l Backend) GetChannelsBalance(ctx context.Context) (*models.ChannelsBalance, error) {
	resp, err := l.client.Call("listfunds")
	if err != nil {
		return nil, err
	}

	balance := models.ChannelsBalance{}
	resp.Get("channels").ForEach(func(_, ch gjson.Result) bool {
		balance.Balance += ch.Get("channel_sat").Int()
		return true
	})

	return &balance, nil
}

func (l Backend) ListChannels(ctx context.Context, opt ...options.Channel) ([]*models.Channel, error) {
	resp, err := l.client.Call("listpeers")
	if err != nil {
		return nil, err
	}

	channels := make([]*models.Channel, 0, resp.Get("peers.#").Int())
	i := 0
	resp.Get("peers").ForEach(func(_, peer gjson.Result) bool {
		if peer.Get("channels.#").Int() == 0 {
			return true
		}
		ch := peer.Get("channels.0")

		channels = append(channels, &models.Channel{
			ID:                  uint64(i),
			Status:              models.ChannelActive,
			RemotePubKey:        peer.Get("id").String(),
			ChannelPoint:        ch.Get("short_channel_id").String(),
			Capacity:            ch.Get("msatoshi_total").Int(),
			LocalBalance:        ch.Get("msatoshi_to_us").Int(),
			RemoteBalance:       (ch.Get("msatoshi_total").Int() - ch.Get("msatoshi_to_us").Int()),
			CommitFee:           0,
			CommitWeight:        0,
			FeePerKiloWeight:    0,
			UnsettledBalance:    0,
			TotalAmountSent:     ch.Get("in_msatoshi_fulfilled").Int(),
			TotalAmountReceived: ch.Get("out_msatoshi_fulfilled").Int(),
			ConfirmationHeight:  nil,
			UpdatesCount:        uint64(ch.Get("in_payments_fulfilled").Int() + ch.Get("out_payments_fulfilled").Int()),
			CSVDelay:            uint32(ch.Get("our_to_self_delay").Int()),
			Private:             ch.Get("private").Bool(),
			PendingHTLC:         []*models.HTLC{},
			LastUpdate:          nil,
			Node:                nil,
			Policy1: &models.RoutingPolicy{
				MinHtlc: ch.Get("htlc_minimum_msat").Int(),
			},
			Policy2: &models.RoutingPolicy{
				MinHtlc: ch.Get("htlc_minimum_msat").Int(),
			},
		})

		i++
		return true
	})

	return channels, nil
}

func (l Backend) GetChannelInfo(ctx context.Context, channel *models.Channel) error {
	res, err := l.client.Call("listchannels", channel.ChannelPoint)
	if err != nil {
		return err
	}

	// the channel is represented as if it was two
	fromUs := res.Get(`channels.#[destination=="` + channel.RemotePubKey + `"]`)
	toUs := res.Get(`channels.#[source=="` + channel.RemotePubKey + `"]`)

	// last update
	lastUpdate := fromUs.Get("last_update").Int()
	if toUs.Get("last_update").Int() > lastUpdate {
		lastUpdate = toUs.Get("last_update").Int()
	}
	t := time.Unix(lastUpdate, 0)
	channel.LastUpdate = &t

	// routing policies
	channel.Policy1.TimeLockDelta = uint32(fromUs.Get("delay").Int())
	channel.Policy1.FeeBaseMsat = fromUs.Get("base_fee_millisatoshi").Int()
	channel.Policy1.FeeRateMilliMsat = fromUs.Get("fee_per_millionth").Int()
	channel.Policy1.Disabled = !fromUs.Get("active").Bool()
	channel.Policy2.TimeLockDelta = uint32(toUs.Get("delay").Int())
	channel.Policy2.FeeBaseMsat = toUs.Get("base_fee_millisatoshi").Int()
	channel.Policy2.FeeRateMilliMsat = toUs.Get("fee_per_millionth").Int()
	channel.Policy2.Disabled = !toUs.Get("active").Bool()

	return nil
}

func (l Backend) GetNode(ctx context.Context, pubkey string) (*models.Node, error) {
	nodeinfo := models.Node{
		PubKey: pubkey,
	}

	if res, err := l.client.Call("listnodes", pubkey); err == nil {
		node := res.Get("nodes.0")
		nodeinfo.Alias = node.Get("alias").String()
		addresses := make([]*models.NodeAddress, node.Get("addresses.#").Int())
		i := 0
		node.Get("addresses").ForEach(func(_, addr gjson.Result) bool {
			addresses[i] = &models.NodeAddress{
				Network: addr.Get("type").String(),
				Addr:    addr.Get("addr").String() + ":" + addr.Get("port").String(),
			}
			i++
			return true
		})
		nodeinfo.Addresses = addresses
	} else {
		return nil, err
	}

	if res, err := l.client.CallNamed("listchannels", "source", pubkey); err == nil {
		nodeinfo.NumChannels = uint32(res.Get("channels.#").Int())

		var lastUpdate int64
		res.Get("channels").ForEach(func(_, ch gjson.Result) bool {
			nodeinfo.TotalCapacity += ch.Get("satoshis").Int()

			thisLastUpdate := ch.Get("last_update").Int()
			if lastUpdate < thisLastUpdate {
				lastUpdate = thisLastUpdate
			}
			return true
		})
		t := time.Unix(lastUpdate, 0)
		nodeinfo.LastUpdate = t
	} else {
		return nil, err
	}

	return &nodeinfo, nil
}

func (l Backend) CreateInvoice(ctx context.Context, amount int64, desc string) (*models.Invoice, error) {
	return nil, nil
}

func (l Backend) GetInvoice(ctx context.Context, RHash string) (*models.Invoice, error) {
	return nil, nil
}

func (l Backend) SendPayment(ctx context.Context, payreq *models.PayReq) (*models.Payment, error) {
	return nil, nil
}

func (l Backend) DecodePayReq(ctx context.Context, payreq string) (*models.PayReq, error) {
	return nil, nil
}
