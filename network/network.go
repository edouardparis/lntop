package network

import (
	"context"

	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network/lnd"
	"github.com/edouardparis/lntop/network/mock"
	"github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/network/options"
)

type Network interface {
	Ping() error

	SubscribeInvoice(context.Context, chan *models.Invoice) error

	SubscribeChannels(context.Context, chan *models.ChannelUpdate) error

	NodeName() string

	Info(context.Context) (*models.Info, error)

	GetNode(context.Context, string) (*models.Node, error)

	GetWalletBalance(context.Context) (*models.WalletBalance, error)

	GetChannelsBalance(context.Context) (*models.ChannelsBalance, error)

	ListChannels(context.Context, ...options.Channel) ([]*models.Channel, error)

	GetChannelInfo(context.Context, *models.Channel) error

	CreateInvoice(context.Context, int64, string) (*models.Invoice, error)

	GetInvoice(context.Context, string) (*models.Invoice, error)

	DecodePayReq(context.Context, string) (*models.PayReq, error)

	SendPayment(context.Context, *models.PayReq) (*models.Payment, error)

	GetTransactions(context.Context) ([]*models.Transaction, error)

	SubscribeTransactions(context.Context, chan *models.Transaction) error
}

func New(c *config.Network, logger logging.Logger) (Network, error) {
	var (
		err error
		b   Network
	)
	if c.Type == "mock" {
		b = mock.New(c)
	} else {
		b, err = lnd.New(c, logger.With(logging.String("network", "lnd")))
		if err != nil {
			return nil, err
		}
	}

	err = b.Ping()
	if err != nil {
		return nil, err
	}

	return b, nil
}
