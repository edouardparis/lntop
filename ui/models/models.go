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
	WalletBalance   *WalletBalance
	ChannelsBalance *ChannelsBalance
}

func New(app *app.App) *Models {
	return &Models{
		App:             app,
		Info:            &Info{},
		Channels:        &Channels{},
		WalletBalance:   &WalletBalance{},
		ChannelsBalance: &ChannelsBalance{},
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

type Channels struct {
	Items []*models.Channel
}

func (m *Models) RefreshChannels(ctx context.Context) error {
	channels, err := m.App.Network.ListChannels(ctx)
	if err != nil {
		return err
	}
	*m.Channels = Channels{channels}
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
