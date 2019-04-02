package ui

import (
	"context"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/events"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/ui/models"
	"github.com/edouardparis/lntop/ui/views"
)

type controller struct {
	logger logging.Logger
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

func (c *controller) SetModels(ctx context.Context) error {
	err := c.models.RefreshInfo(ctx)
	if err != nil {
		return err
	}

	err = c.models.RefreshWalletBalance(ctx)
	if err != nil {
		return err
	}

	err = c.models.RefreshChannelsBalance(ctx)
	if err != nil {
		return err
	}

	return c.models.RefreshChannels(ctx)
}

func (c *controller) Listen(ctx context.Context, g *gocui.Gui, sub chan *events.Event) {
	c.logger.Debug("Listening...")
	for event := range sub {
		var err error
		switch event.Type {
		case events.BlockReceived:
			err = c.models.RefreshInfo(ctx)
		case events.ChannelPending:
			err = c.models.RefreshInfo(ctx)
		case events.ChannelActive:
			err = c.models.RefreshInfo(ctx)
		case events.ChannelInactive:
			err = c.models.RefreshInfo(ctx)
		case events.PeerUpdated:
			err = c.models.RefreshInfo(ctx)
		default:
			c.logger.Info("event received", logging.String("type", event.Type))
		}
		if err != nil {
			c.logger.Error("failed", logging.String("event", event.Type))
		}
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *controller) Help(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	view := c.views.Get(g.CurrentView().Name())
	if view == nil {
		return nil
	}

	if view.Name() != views.HELP {
		c.views.SetPrevious(view)
		return c.views.Help.Set(g, 0, -1, maxX, maxY)
	}

	err := g.DeleteView(views.HELP)
	if err != nil {
		return err
	}

	if c.views.Previous != nil {
		_, err := g.SetCurrentView(c.views.Previous.Name())
		return err
	}

	return nil
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

	err = g.SetKeybinding("", gocui.KeyF1, gocui.ModNone, c.Help)
	if err != nil {
		return err
	}

	return nil
}

func newController(app *app.App) *controller {
	m := models.New(app)
	return &controller{
		logger: app.Logger.With(logging.String("logger", "controller")),
		models: m,
		views:  views.New(m),
	}
}
