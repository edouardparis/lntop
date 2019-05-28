package ui

import (
	"context"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/events"
)

type Pubsub interface {
	Subscribe(events.Publisher)
	Unsubscribe(events.Publisher)
	Events() chan events.Event
	Stop() error
}

func Run(ctx context.Context, app *app.App, ps Pubsub) error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Cursor = true
	ctrl := newController(app)
	err = ctrl.SetModels(ctx)
	if err != nil {
		return err
	}

	g.SetManagerFunc(ctrl.layout)

	err = setKeyBinding(ctrl, g)
	if err != nil {
		return err
	}

	go ctrl.Listen(ctx, g, ps.Events())

	err = g.MainLoop()
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(ps.Stop())
}
