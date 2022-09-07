package mock

import (
	"context"
	"crypto/sha256"
	"fmt"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/network/options"
)

type Backend struct {
	invoices map[string]models.Invoice
	count    uint64
	cfg      *config.Network
	sync.RWMutex
}

func (b *Backend) Ping() error {
	return nil
}

func (b *Backend) Info(ctx context.Context) (*models.Info, error) {
	return nil, nil
}

func (l *Backend) SendPayment(ctx context.Context, payreq *models.PayReq) (*models.Payment, error) {
	return nil, nil
}

func (b *Backend) NodeName() string {
	return b.cfg.Name
}

func (b *Backend) SubscribeInvoice(ctx context.Context, ChannelInvoice chan *models.Invoice) error {
	return nil
}

func (b *Backend) SubscribeChannels(context.Context, chan *models.ChannelUpdate) error {
	return nil
}

func (b *Backend) SubscribeTransactions(ctx context.Context, channel chan *models.Transaction) error {
	return nil
}

func (b *Backend) SubscribeRoutingEvents(ctx context.Context, channel chan *models.RoutingEvent) error {
	return nil
}

func (b *Backend) GetNode(ctx context.Context, pubkey string, includeChannels bool) (*models.Node, error) {
	return &models.Node{}, nil
}

func (b *Backend) GetWalletBalance(ctx context.Context) (*models.WalletBalance, error) {
	return &models.WalletBalance{}, nil
}

func (b *Backend) GetTransactions(ctx context.Context) ([]*models.Transaction, error) {
	return []*models.Transaction{}, nil
}

func (b *Backend) GetChannelsBalance(ctx context.Context) (*models.ChannelsBalance, error) {
	return &models.ChannelsBalance{}, nil
}

func (b *Backend) ListChannels(ctx context.Context, opt ...options.Channel) ([]*models.Channel, error) {
	return []*models.Channel{}, nil
}

func (b *Backend) GetChannelInfo(ctx context.Context, channel *models.Channel) error {
	return nil
}

func (b *Backend) DecodePayReq(ctx context.Context, payreq string) (*models.PayReq, error) {
	return &models.PayReq{}, nil
}

func (b *Backend) CreateInvoice(ctx context.Context, amt int64, desc string) (*models.Invoice, error) {
	b.Lock()
	defer b.Unlock()
	b.count++

	key := uuid.Must(uuid.NewV4()).String()

	preimage := []byte(fmt.Sprintf("preimage %s", key))
	hash := sha256.Sum256([]byte(preimage))

	invoice := &models.Invoice{
		Index:          b.count,
		RPreImage:      preimage,
		RHash:          hash[:],
		Amount:         amt,
		Description:    desc,
		CreationDate:   time.Now().Unix(),
		Expiry:         3600,
		PaymentRequest: "lnbc28600u1pw9n7g7pp5enjn8exsyymyl6mlxmcvy7fdcwuh04z96swfmtasznppglgdyvsqdqqcqzysc8rve6vdwuvketcn7yp8gu3ltvq29vj588erp3at9z2msqj0yhhjdwsf7qtfy5lwf8favm6u3wr5qklvprlhrz89pknpdfxnc55wy6sqnrxjh7",
	}

	b.invoices[string(invoice.RHash)] = *invoice

	return invoice, nil
}

func (b *Backend) GetInvoice(ctx context.Context, hash string) (*models.Invoice, error) {
	invoice, ok := b.invoices[hash]
	if !ok {
		return nil, errors.New("unable to locate invoice")
	}

	return &invoice, nil
}

func New(c *config.Network) *Backend {
	return &Backend{
		invoices: make(map[string]models.Invoice),
		cfg:      c,
	}
}
