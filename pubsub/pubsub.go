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
	logger  logging.Logger
	network *network.Network
	wg      *sync.WaitGroup
}

func New(logger logging.Logger, network *network.Network) *PubSub {
	return &PubSub{
		logger:  logger.With(logging.String("logger", "pubsub")),
		network: network,
		wg:      &sync.WaitGroup{},
		stop:    make(chan bool),
	}
}

func (p *PubSub) invoices(ctx context.Context, sub chan *events.Event) {
	p.wg.Add(3)
	invoices := make(chan *models.Invoice)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for invoice := range invoices {
			p.logger.Info("receive invoice", logging.Object("invoice", invoice))
			if invoice.Settled {
				sub <- events.New(events.InvoiceSettled)
			} else {
				sub <- events.New(events.InvoiceCreated)
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

func (p *PubSub) transactions(ctx context.Context, sub chan *events.Event) {
	p.wg.Add(3)
	transactions := make(chan *models.Transaction)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for tx := range transactions {
			p.logger.Info("receive transaction", logging.String("tx_hash", tx.TxHash))
			sub <- events.New(events.TransactionCreated)
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

func (p *PubSub) routingUpdates(ctx context.Context, sub chan *events.Event) {
	p.wg.Add(3)
	routingUpdates := make(chan *models.RoutingEvent)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for hu := range routingUpdates {
			p.logger.Info("receive htlcUpdate")
			sub <- events.NewWithData(events.RoutingEventUpdated, hu)
		}
		p.wg.Done()
	}()

	go func() {
		err := p.network.SubscribeRoutingEvents(ctx, routingUpdates)
		if err != nil {
			p.logger.Error("SubscribeRoutingEvents returned an error", logging.Error(err))
		}
		p.wg.Done()
	}()

	go func() {
		<-p.stop
		cancel()
		close(routingUpdates)
		p.wg.Done()
	}()
}

func (p *PubSub) Stop() {
	p.stop <- true
	close(p.stop)
	p.logger.Debug("Received signal, gracefully stopping")
}

func (p *PubSub) Run(ctx context.Context, sub chan *events.Event) {
	p.logger.Debug("Starting...")

	p.invoices(ctx, sub)
	p.transactions(ctx, sub)
	p.routingUpdates(ctx, sub)
	p.ticker(ctx, sub,
		withTickerInfo(),
		withTickerChannelsBalance(),
		// no need for ticker Wallet balance, transactions subscriber is enough
		// withTickerWalletBalance(),
	)

	<-p.stop
	p.wg.Wait()
}
