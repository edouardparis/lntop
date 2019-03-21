package ui

import (
	"context"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/ui/views"
)

type Ui struct {
	app      *app.App
	channels *views.Channels
}

func (u *Ui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.Cursor = true
	g.SetManagerFunc(u.layout)

	err = u.setKeyBinding(g)
	if err != nil {
		return err
	}

	g.Update(u.Refresh)

	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		return err
	}

	return err
}

func (u *Ui) setKeyBinding(g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, cursorUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, cursorDown)
	if err != nil {
		return err
	}

	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *Ui) Refresh(g *gocui.Gui) error {
	channels, err := u.app.Network.ListChannels(context.Background())
	if err != nil {
		return err
	}
	u.channels.Update(channels)
	return nil
}

func (u *Ui) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	return u.channels.Set(g, 0, maxY/8, maxX-1, maxY-1)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func New(app *app.App) *Ui {
	return &Ui{
		app:      app,
		channels: views.NewChannels(),
	}
}
