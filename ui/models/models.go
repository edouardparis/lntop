package models

import (
	"context"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/network/models"
)

type Models struct {
	App  *app.App
	Info *Info
}

func New(app *app.App) *Models {
	return &Models{
		App:  app,
		Info: &Info{},
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

type Info struct {
	*models.Info
}
