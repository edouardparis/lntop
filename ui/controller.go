package ui

import (
	"context"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/ui/models"
	"github.com/edouardparis/lntop/ui/views"
)

type controller struct {
	models *models.Models
	views  *views.Views
}

func (c *controller) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	return c.views.Layout(g, maxX, maxY)
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

func (c *controller) Update(ctx context.Context) error {
	info, err := c.models.App.Network.Info(ctx)
	if err != nil {
		return err
	}
	alias := info.Alias
	if c.models.App.Config.Network.Name != "" {
		alias = c.models.App.Config.Network.Name
	}
	c.views.Header.Update(alias, "lnd", info.Version)
	c.views.Summary.UpdateChannelsStats(
		info.NumPendingChannels,
		info.NumActiveChannels,
		info.NumInactiveChannels,
	)

	channels, err := c.models.App.Network.ListChannels(ctx)
	if err != nil {
		return err
	}
	c.views.Channels.Update(channels)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *controller) setKeyBinding(g *gocui.Gui) error {
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

func newController(app *app.App) *controller {
	return &controller{
		models: models.New(app),
		views:  views.New(),
	}
}
