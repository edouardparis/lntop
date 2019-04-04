package backend

import (
	"context"

	"github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/network/options"
)

type Backend interface {
	SubscribeInvoice(context.Context, chan *models.Invoice) error

	SubscribeChannels(context.Context, chan *models.ChannelUpdate) error

	NodeName() string

	Info(ctx context.Context) (*models.Info, error)

	GetWalletBalance(context.Context) (*models.WalletBalance, error)

	GetChannelsBalance(context.Context) (*models.ChannelsBalance, error)

	ListChannels(context.Context, ...options.Channel) ([]*models.Channel, error)

	GetChannelInfo(context.Context, *models.Channel) error

	CreateInvoice(context.Context, int64, string) (*models.Invoice, error)

	GetInvoice(context.Context, string) (*models.Invoice, error)

	DecodePayReq(context.Context, string) (*models.PayReq, error)

	SendPayment(context.Context, *models.PayReq) (*models.Payment, error)
}
