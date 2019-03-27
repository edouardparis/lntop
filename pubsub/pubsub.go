package pubsub

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/events"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network"
	"github.com/edouardparis/lntop/network/models"
)

type pubSub struct {
	stop    chan bool
	logger  logging.Logger
	network *network.Network
	wg      *sync.WaitGroup
}

func newPubSub(logger logging.Logger, network *network.Network) *pubSub {
	return &pubSub{
		logger:  logger.With(logging.String("logger", "pubsub")),
		network: network,
		wg:      &sync.WaitGroup{},
		stop:    make(chan bool),
	}
}

func (p *pubSub) invoices(ctx context.Context, sub chan *events.Event) {
	p.wg.Add(2)
	invoices := make(chan *models.Invoice)

	go func() {
		for {
			select {
			case <-p.stop:
				p.wg.Done()
				return
			default:
				invoice := <-invoices
				p.logger.Debug("receive invoice", logging.Object("invoice", invoice))
				if invoice.Settled {
					sub <- events.New(events.InvoiceSettled)
				} else {
					sub <- events.New(events.InvoiceCreated)
				}
			}
		}
	}()

	go func() {
		select {
		case <-p.stop:
			p.wg.Done()
			return
		default:
			err := p.network.SubscribeInvoice(ctx, invoices)
			if err != nil {
				p.logger.Error("SubscribeInvoice returned an error", logging.Error(err))
			}
		}
	}()
}

func (p *pubSub) wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	p.wg.Add(1)
	go func() {
		for sig := range c {
			p.logger.Debug("Received signal, gracefully stopping", logging.String("sig", sig.String()))
			p.wg.Done()
			close(p.stop)
		}
	}()
	p.wg.Wait()
}

func Run(ctx context.Context, app *app.App, sub chan *events.Event) error {
	pubSub := newPubSub(app.Logger, app.Network)
	pubSub.logger.Debug("Starting...")

	pubSub.invoices(ctx, sub)
	pubSub.wait()
	return nil
}