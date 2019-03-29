package ui

import (
	"context"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/events"
)

func Run(ctx context.Context, app *app.App, sub chan *events.Event) error {
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

	err = ctrl.setKeyBinding(g)
	if err != nil {
		return err
	}

	go func() {
		err := ctrl.Refresh(ctx, sub)
		if err != nil {
			g.Update(func(*gocui.Gui) error { return err })
		}
	}()

	err = g.MainLoop()
	close(sub)

	return err
}
