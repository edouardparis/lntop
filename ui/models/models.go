package models

import (
	"context"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/network/models"
)

type Models struct {
	App             *app.App
	Info            *Info
	Channels        *Channels
	CurrentChannel  *Channel
	WalletBalance   *WalletBalance
	ChannelsBalance *ChannelsBalance
}

func New(app *app.App) *Models {
	return &Models{
		App:             app,
		Info:            &Info{},
		Channels:        NewChannels(),
		WalletBalance:   &WalletBalance{},
		ChannelsBalance: &ChannelsBalance{},
		CurrentChannel:  &Channel{},
	}
}

type Info struct {
	*models.Info
}

func (m *Models) RefreshInfo(ctx context.Context) error {
	info, err := m.App.Network.Info(ctx)
	if err != nil {
		return err
	}
	*m.Info = Info{info}
	return nil
}

func (m *Models) RefreshChannels(ctx context.Context) error {
	channels, err := m.App.Network.ListChannels(ctx)
	if err != nil {
		return err
	}
	for i := range channels {
		if !m.Channels.Contains(channels[i]) {
			m.Channels.Add(channels[i])
		}
		channel := m.Channels.GetByID(channels[i].ID)
		if channel != nil &&
			(channel.UpdatesCount < channels[i].UpdatesCount ||
				channel.LastUpdate == nil) {
			err := m.App.Network.GetChannelInfo(ctx, channels[i])
			if err != nil {
				return err
			}

			if channels[i].Node == nil {
				channels[i].Node, err = m.App.Network.GetNode(ctx,
					channels[i].RemotePubKey)
				if err != nil {
					return err
				}
			}

			m.Channels.Update(channels[i])
		}
	}
	return nil
}

func (m *Models) SetCurrentChannel(ctx context.Context, index int) error {
	channel := m.Channels.Get(index)
	if channel == nil {
		return nil
	}
	*m.CurrentChannel = Channel{Item: channel}
	return nil
}

type WalletBalance struct {
	*models.WalletBalance
}

func (m *Models) RefreshWalletBalance(ctx context.Context) error {
	balance, err := m.App.Network.GetWalletBalance(ctx)
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
	balance, err := m.App.Network.GetChannelsBalance(ctx)
	if err != nil {
		return err
	}
	*m.ChannelsBalance = ChannelsBalance{balance}
	return nil
}
