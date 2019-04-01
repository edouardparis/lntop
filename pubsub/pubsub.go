package pubsub

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

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

func (p *pubSub) ticker(ctx context.Context, sub chan *events.Event) {
	p.wg.Add(1)
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		var old *models.Info
		for {
			select {
			case <-p.stop:
				ticker.Stop()
				p.wg.Done()
				return
			case <-ticker.C:
				info, err := p.network.Info(ctx)
				if err != nil {
					p.logger.Error("SubscribeInvoice returned an error", logging.Error(err))
				}
				if old != nil {
					if old.BlockHeight != info.BlockHeight {
						sub <- events.New(events.BlockReceived)
					}

					if old.NumPeers != info.NumPeers {
						sub <- events.New(events.PeerUpdated)
					}

					if old.NumPendingChannels < info.NumPendingChannels {
						sub <- events.New(events.ChannelPending)
					}

					if old.NumActiveChannels < info.NumActiveChannels {
						sub <- events.New(events.ChannelActive)
					}

					if old.NumInactiveChannels < info.NumInactiveChannels {
						sub <- events.New(events.ChannelInactive)
					}
				}
				old = info
			}
		}
	}()
}

func (p *pubSub) invoices(ctx context.Context, sub chan *events.Event) {
	p.wg.Add(3)
	invoices := make(chan *models.Invoice)
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for invoice := range invoices {
			p.logger.Debug("receive invoice", logging.Object("invoice", invoice))
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

func (p *pubSub) wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	p.wg.Add(1)
	go func() {
		sig := <-c
		p.logger.Debug("Received signal, gracefully stopping", logging.String("sig", sig.String()))
		p.stop <- true
		close(p.stop)
		p.wg.Done()
	}()
	p.wg.Wait()
}

func Run(ctx context.Context, app *app.App, sub chan *events.Event) {
	pubSub := newPubSub(app.Logger, app.Network)
	pubSub.logger.Debug("Starting...")

	pubSub.invoices(ctx, sub)
	pubSub.ticker(ctx, sub)
	pubSub.wait()
}
