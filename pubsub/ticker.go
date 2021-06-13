package pubsub

import (
	"context"
	"time"

	"github.com/edouardparis/lntop/events"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network"
	"github.com/edouardparis/lntop/network/models"
)

type tickerFunc func(context.Context, logging.Logger, *network.Network, chan *events.Event)

func (p *PubSub) ticker(ctx context.Context, sub chan *events.Event, fn ...tickerFunc) {
	p.wg.Add(1)
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-p.stop:
				ticker.Stop()
				p.wg.Done()
				return
			case <-ticker.C:
				for i := range fn {
					fn[i](ctx, p.logger, p.network, sub)
				}
			}
		}
	}()
}

// withTickerInfo checks if general information did not changed changed in the ticker interval.
func withTickerInfo() tickerFunc {
	var old *models.Info
	return func(ctx context.Context, logger logging.Logger, net *network.Network, sub chan *events.Event) {
		info, err := net.Info(ctx)
		if err != nil {
			logger.Error("network info returned an error", logging.Error(err))
		}
		if old != nil && info != nil {
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

// withTickerChannelsBalance checks if channels balance and pending balance
// changed in the ticker interval.
func withTickerChannelsBalance() tickerFunc {
	var old *models.ChannelsBalance
	return func(ctx context.Context, logger logging.Logger, net *network.Network, sub chan *events.Event) {
		channelsBalance, err := net.GetChannelsBalance(ctx)
		if err != nil {
			logger.Error("network channels balance returned an error", logging.Error(err))
		}
		if old != nil && channelsBalance != nil {
			if old.Balance != channelsBalance.Balance ||
				old.PendingOpenBalance != channelsBalance.PendingOpenBalance {
				sub <- events.New(events.ChannelBalanceUpdated)
			}
		}
		old = channelsBalance
	}
}

// withTickerWalletBalance checks if wallet balance and pending balance
// changed in the ticker interval.
func withTickerWalletBalance() tickerFunc {
	var old *models.WalletBalance
	return func(ctx context.Context, logger logging.Logger, net *network.Network, sub chan *events.Event) {
		walletBalance, err := net.GetWalletBalance(ctx)
		if err != nil {
			logger.Error("network wallet balance returned an error", logging.Error(err))
		}
		if old != nil && walletBalance != nil {
			if old.TotalBalance != walletBalance.TotalBalance ||
				old.ConfirmedBalance != walletBalance.ConfirmedBalance ||
				old.UnconfirmedBalance != walletBalance.UnconfirmedBalance {
				sub <- events.New(events.WalletBalanceUpdated)
			}
		}
		old = walletBalance
	}
}
