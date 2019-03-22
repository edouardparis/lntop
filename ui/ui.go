package ui

import (
	"context"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/app"
)

func Run(ctx context.Context, app *app.App) error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Cursor = true
	ctrl := newController(app)
	g.SetManagerFunc(ctrl.layout)

	err = ctrl.setKeyBinding(g)
	if err != nil {
		return err
	}

	err = ctrl.Update(ctx)
	if err != nil {
		return err
	}

	g.Update(ctrl.Refresh(ctx))

	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		return err
	}

	return err
}
