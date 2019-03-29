package models

import (
	"context"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/network/models"
)

type Models struct {
	App      *app.App
	Info     *Info
	Channels *Channels
}

func New(app *app.App) *Models {
	return &Models{
		App:      app,
		Info:     &Info{},
		Channels: &Channels{},
	}
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
	*m.Channels = Channels{channels}
	return nil
}

type Info struct {
	*models.Info
}

type Channels struct {
	Items []*models.Channel
}
