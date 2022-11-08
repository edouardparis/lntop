package models

import (
	"context"
	"strconv"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network"
	"github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/network/options"
)

type Models struct {
	logger          logging.Logger
	network         *network.Network
	Info            *Info
	Channels        *Channels
	WalletBalance   *WalletBalance
	ChannelsBalance *ChannelsBalance
	Transactions    *Transactions
	RoutingLog      *RoutingLog
	FwdingHist      *FwdingHist
}

func New(app *app.App) *Models {
	fwdingHist := FwdingHist{}
	startTime := app.Config.Views.FwdingHist.Options.GetOption("START_TIME", "start_time")
	maxNumEvents := app.Config.Views.FwdingHist.Options.GetOption("MAX_NUM_EVENTS", "max_num_events")

	if startTime != "" {
		fwdingHist.StartTime = startTime
	}

	if maxNumEvents != "" {
		max, err := strconv.ParseUint(maxNumEvents, 10, 32)
		if err != nil {
			app.Logger.Info("Couldn't parse the maximum number of forwarding events.")
		} else {
			fwdingHist.MaxNumEvents = uint32(max)
		}
	}

	return &Models{
		logger:          app.Logger.With(logging.String("logger", "models")),
		network:         app.Network,
		Info:            &Info{},
		Channels:        NewChannels(),
		WalletBalance:   &WalletBalance{},
		ChannelsBalance: &ChannelsBalance{},
		Transactions:    &Transactions{},
		RoutingLog:      &RoutingLog{},
		FwdingHist:      &fwdingHist,
	}
}

type Info struct {
	*models.Info
}

func (m *Models) RefreshInfo(ctx context.Context) error {
	info, err := m.network.Info(ctx)
	if err != nil {
		return err
	}
	*m.Info = Info{info}
	return nil
}

func (m *Models) RefreshForwardingHistory(ctx context.Context) error {
<<<<<<< Updated upstream
	forwardingEvents, err := m.network.GetForwardingHistory(ctx)
=======
	forwardingEvents, err := m.network.GetForwardingHistory(ctx, m.FwdingHist.StartTime, m.FwdingHist.MaxNumEvents)
>>>>>>> Stashed changes
	if err != nil {
		return err
	}

	m.FwdingHist.Update(forwardingEvents)

	return nil
}

func (m *Models) RefreshChannels(ctx context.Context) error {
	channels, err := m.network.ListChannels(ctx, options.WithChannelPending)
	if err != nil {
		return err
	}
	index := map[string]*models.Channel{}
	for i := range channels {
		index[channels[i].ChannelPoint] = channels[i]
		if channels[i].ID > 0 {
			channels[i].Age = m.Info.BlockHeight - uint32(channels[i].ID>>40)
		}
		if !m.Channels.Contains(channels[i]) {
			m.Channels.Add(channels[i])
		}
		channel := m.Channels.GetByChanPoint(channels[i].ChannelPoint)
		if channel != nil &&
			(channel.UpdatesCount < channels[i].UpdatesCount ||
				channel.LastUpdate == nil || channel.LocalPolicy == nil || channel.RemotePolicy == nil) {
			err := m.network.GetChannelInfo(ctx, channels[i])
			if err != nil {
				return err
			}

			if channels[i].Node == nil {
				channels[i].Node, err = m.network.GetNode(ctx,
					channels[i].RemotePubKey, false)
				if err != nil {
					m.logger.Debug("refreshChannels: cannot find Node",
						logging.String("pubkey", channels[i].RemotePubKey))
				}
			}
		}

		m.Channels.Update(channels[i])
	}
	for _, c := range m.Channels.List() {
		if _, ok := index[c.ChannelPoint]; !ok {
			c.Status = models.ChannelClosed
		}
	}
	return nil
}

type WalletBalance struct {
	*models.WalletBalance
}

func (m *Models) RefreshWalletBalance(ctx context.Context) error {
	balance, err := m.network.GetWalletBalance(ctx)
	if err != nil {
		return err
	}
	*m.WalletBalance = WalletBalance{balance}
	return nil
}

type ChannelsBalance struct {
	*models.ChannelsBalance
}

func (m *Models) RefreshChannelsBalance(ctx context.Context) error {
	balance, err := m.network.GetChannelsBalance(ctx)
	if err != nil {
		return err
	}
	*m.ChannelsBalance = ChannelsBalance{balance}
	return nil
}

type RoutingLog struct {
	Log []*models.RoutingEvent
}

const MaxRoutingEvents = 512 // 8K monitor @ 8px per line = 540

func (m *Models) RefreshRouting(update interface{}) func(context.Context) error {
	return (func(ctx context.Context) error {
		hu, ok := update.(*models.RoutingEvent)
		if ok {
			found := false
			for _, hlu := range m.RoutingLog.Log {
				if hlu.Equals(hu) {
					hlu.Update(hu)
					found = true
					break
				}
			}
			if !found {
				if len(m.RoutingLog.Log) == MaxRoutingEvents {
					m.RoutingLog.Log = m.RoutingLog.Log[1:]
				}
				m.RoutingLog.Log = append(m.RoutingLog.Log, hu)
			}
		} else {
			m.logger.Error("refreshRouting: invalid event data")
		}
		return nil
	})
}

func (m *Models) RefreshPolicies(update interface{}) func(context.Context) error {
	return func(ctx context.Context) error {
		for _, chanpoint := range update.(*models.ChannelEdgeUpdate).ChanPoints {
			if m.Channels.Contains(&models.Channel{ChannelPoint: chanpoint}) {
				m.logger.Debug("updating channel", logging.String("chanpoint", chanpoint))
				channel := m.Channels.GetByChanPoint(chanpoint)
				err := m.network.GetChannelInfo(ctx, channel)
				if err != nil {
					m.logger.Error("error updating channel info", logging.Error(err))
				}
			}
		}
		return nil
	}
}

func (m *Models) RefreshCurrentNode(ctx context.Context) (err error) {
	cur := m.Channels.Current()
	if cur != nil {
		m.Channels.CurrentNode, err = m.network.GetNode(ctx, cur.RemotePubKey, true)
	}
	return
}
