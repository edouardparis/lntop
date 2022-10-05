package ui

import (
	"context"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/events"
)

func Run(ctx context.Context, app *app.App, sub chan *events.Event) error {
	g, err := gocui.NewGui(gocui.Output256)
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

	go ctrl.Listen(ctx, g, sub)

	err = g.MainLoop()

	return errors.WithStack(err)
}
