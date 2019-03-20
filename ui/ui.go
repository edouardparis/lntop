package ui

import (
	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/ui/views"
)

type Ui struct {
	gui      *gocui.Gui
	channels *views.Channels
}

func (u *Ui) Run() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	u.gui = g

	g.Cursor = true
	g.SetManagerFunc(u.layout)

	err = g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		return err
	}

	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		return err
	}

	return err
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
		channels: views.NewChannels(app.Network),
	}
}
