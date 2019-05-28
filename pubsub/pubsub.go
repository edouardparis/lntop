package pubsub

import (
	"context"
	"sync"

	"github.com/edouardparis/lntop/events"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network"
	"github.com/edouardparis/lntop/network/models"
)

type PubSub struct {
	stop    chan bool
	sub     chan events.Event
	logger  logging.Logger
	network network.Network
	wg      *sync.WaitGroup
}

func New(logger logging.Logger, network network.Network) *PubSub {
	return &PubSub{
		logger:  logger.With(logging.String("logger", "pubsub")),
		network: network,
		wg:      &sync.WaitGroup{},
		stop:    make(chan bool),
		sub:     make(chan events.Event),
	}
}

func (p *PubSub) invoices(ctx context.Context) {
	p.wg.Add(3)
	invoices := make(chan *models.Invoice)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for invoice := range invoices {
			p.logger.Debug("receive invoice", logging.Object("invoice", invoice))
			if invoice.Settled {
				p.sub <- events.InvoiceSettled
			} else {
				p.sub <- events.InvoiceCreated
			}
		}
		p.wg.Done()
	}()

	go func() {
		err := p.network.SubscribeInvoice(ctx, invoices)
		if err != nil {
			p.logger.Error("SubscribeInvoice returned an error", logging.Error(err))
		}
		p.wg.Done()
	}()

	go func() {
		<-p.stop
		cancel()
		close(invoices)
		p.wg.Done()
	}()
}

func (p *PubSub) transactions(ctx context.Context) {
	p.wg.Add(3)
	transactions := make(chan *models.Transaction)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for tx := range transactions {
			p.logger.Debug("receive transaction", logging.String("tx_hash", tx.TxHash))
			p.sub <- events.TransactionCreated
		}
		p.wg.Done()
	}()

	go func() {
		err := p.network.SubscribeTransactions(ctx, transactions)
		if err != nil {
			p.logger.Error("SubscribeTransactions returned an error", logging.Error(err))
		}
		p.wg.Done()
	}()

	go func() {
		<-p.stop
		cancel()
		close(transactions)
		p.wg.Done()
	}()
}

func (p *PubSub) Subscribe(pub events.Publisher)   {}
func (p *PubSub) Unsubscribe(pub events.Publisher) {}

func (p *PubSub) Events() chan events.Event {
	return p.sub
}

func (p *PubSub) Stop() error {
	p.stop <- true
	close(p.stop)
	close(p.sub)
	p.logger.Debug("Received signal, gracefully stopping")
	return nil
}

func (p *PubSub) Run(ctx context.Context) {
	p.logger.Debug("Starting...")

	p.invoices(ctx)
	p.transactions(ctx)
	p.ticker(ctx, p.sub,
		withTickerInfo(),
		withTickerChannelsBalance(),
		// no need for ticker Wallet balance, transactions subscriber is enough
		// withTickerWalletBalance(),
	)

	<-p.stop
	p.wg.Wait()
}
