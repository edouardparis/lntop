package ui

import (
	"context"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/ui/views"
	"github.com/jroimartin/gocui"
)

type controller struct {
	app      *app.App
	header   *views.Header
	summary  *views.Summary
	channels *views.Channels
}

func (c *controller) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	err := c.header.Set(g, 0, 0, maxX, 2)
	if err != nil {
		return err
	}

	err = c.summary.Set(g, 0, 2, maxX, 7)
	if err != nil {
		return err
	}

	return c.channels.Set(g, 0, 7, maxX-1, maxY-1)
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
	info, err := c.app.Network.Info(ctx)
	if err != nil {
		return err
	}
	c.header.Update(info.Alias, "lnd", info.Version)
	c.summary.UpdateChannelsStats(
		info.NumPendingChannels,
		info.NumActiveChannels,
		info.NumInactiveChannels,
	)

	channels, err := c.app.Network.ListChannels(ctx)
	if err != nil {
		return err
	}
	c.channels.Update(channels)
	return nil
}

func (c *controller) Refresh(ctx context.Context) func(*gocui.Gui) error {
	return func(g *gocui.Gui) error {
		info, err := c.app.Network.Info(ctx)
		if err != nil {
			return err
		}
		c.header.Update(info.Alias, "lnd", info.Version)

		channels, err := c.app.Network.ListChannels(ctx)
		if err != nil {
			return err
		}
		c.channels.Update(channels)
		return nil
	}
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
		app:      app,
		header:   views.NewHeader(),
		summary:  views.NewSummary(),
		channels: views.NewChannels(),
	}
}
