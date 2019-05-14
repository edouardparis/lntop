package ui

import (
	"context"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/events"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/ui/cursor"
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

func (c *controller) cursorDown(g *gocui.Gui, v *gocui.View) error {
	view := c.views.Get(v)
	if view != nil {
		return cursor.Down(view)
	}
	return nil
}

func (c *controller) cursorUp(g *gocui.Gui, v *gocui.View) error {
	view := c.views.Get(v)
	if view != nil {
		return cursor.Up(view)
	}
	return nil
}

func (c *controller) cursorRight(g *gocui.Gui, v *gocui.View) error {
	view := c.views.Get(v)
	if view != nil {
		return cursor.Right(view)
	}
	return nil
}

func (c *controller) cursorLeft(g *gocui.Gui, v *gocui.View) error {
	view := c.views.Get(v)
	if view != nil {
		return cursor.Left(view)
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

	err = c.models.RefreshTransactions(ctx)
	if err != nil {
		return err
	}

	return c.models.RefreshChannels(ctx)
}

func (c *controller) Listen(ctx context.Context, g *gocui.Gui, sub chan *events.Event) {
	c.logger.Debug("Listening...")
	refresh := func(fn ...func(context.Context) error) {
		for i := range fn {
			err := fn[i](ctx)
			if err != nil {
				c.logger.Error("failed", logging.Error(err))
			}
		}
		g.Update(func(*gocui.Gui) error { return nil })
	}

	for event := range sub {
		c.logger.Debug("event received", logging.String("type", event.Type))
		switch event.Type {
		case events.TransactionCreated:
			refresh(c.models.RefreshTransactions)
		case events.BlockReceived:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshTransactions,
			)
		case events.WalletBalanceUpdated:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshWalletBalance,
			)
		case events.ChannelBalanceUpdated:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshChannelsBalance,
				c.models.RefreshChannels,
			)
		case events.ChannelPending:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshChannelsBalance,
				c.models.RefreshChannels,
			)
		case events.ChannelActive:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshChannelsBalance,
				c.models.RefreshChannels,
			)
		case events.ChannelInactive:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshChannelsBalance,
				c.models.RefreshChannels,
			)
		case events.InvoiceSettled:
			refresh(
				c.models.RefreshInfo,
				c.models.RefreshChannelsBalance,
				c.models.RefreshChannels,
			)
		case events.PeerUpdated:
			refresh(c.models.RefreshInfo)
		}
	}
}

func (c *controller) Help(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	view := c.views.Get(g.CurrentView())
	if view == nil {
		return nil
	}

	if view.Name() != views.HELP {
		c.views.Main = view
		return c.views.Help.Set(g, 0, -1, maxX, maxY)
	}

	err := view.Delete(g)
	if err != nil {
		return err
	}

	if c.views.Main != nil {
		_, err := g.SetCurrentView(c.views.Main.Name())
		return err
	}

	return nil
}

func (c *controller) Menu(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if v.Name() == c.views.Help.Name() {
		return nil
	}

	if v.Name() != c.views.Menu.Name() {
		err := c.views.Menu.Set(g, 0, 6, 10, maxY)
		if err != nil {
			return err
		}

		err = c.views.Main.Set(g, 11, 6, maxX-1, maxY)
		if err != nil {
			return err
		}

		_, err = g.SetCurrentView(c.views.Menu.Name())
		return err
	}

	err := c.views.Menu.Delete(g)
	if err != nil {
		return err
	}

	if c.views.Main != nil {
		_, err := g.SetCurrentView(c.views.Main.Name())
		return err
	}

	return nil
}

func (c *controller) Order(order models.Order) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		view := c.views.Get(v)
		if view == nil {
			return nil
		}
		switch view.Name() {
		case views.CHANNELS:
			c.views.Channels.Sort("", order)
		case views.TRANSACTIONS:
			c.views.Transactions.Sort("", order)
		}
		return nil
	}
}

func (c *controller) OnEnter(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	view := c.views.Get(v)
	if view == nil {
		return nil
	}

	switch view.Name() {
	case views.CHANNELS:
		index := c.views.Channels.Index()
		c.models.Channels.SetCurrent(index)
		c.views.Main = c.views.Channel
		return ToggleView(g, view, c.views.Channels)

	case views.CHANNEL:
		c.views.Main = c.views.Channels
		return ToggleView(g, view, c.views.Channels)

	case views.MENU:
		current := c.views.Menu.Current()
		if c.views.Main.Name() == current {
			return nil
		}

		switch current {
		case views.TRANSACTIONS:
			err := c.views.Main.Delete(g)
			if err != nil {
				return err
			}

			c.views.Main = c.views.Transactions
			err = c.views.Transactions.Set(g, 11, 6, maxX-1, maxY)
			if err != nil {
				return err
			}
		case views.CHANNELS:
			err := c.views.Main.Delete(g)
			if err != nil {
				return err
			}

			c.views.Main = c.views.Channels
			err = c.views.Channels.Set(g, 11, 6, maxX-1, maxY)
			if err != nil {
				return err
			}
		}
	case views.TRANSACTIONS:
		index := c.views.Transactions.Index()
		c.models.Transactions.SetCurrent(index)
		c.views.Main = c.views.Transaction
		return ToggleView(g, view, c.views.Transaction)

	case views.TRANSACTION:
		c.views.Main = c.views.Transactions
		return ToggleView(g, view, c.views.Transactions)
	}
	return nil
}

func ToggleView(g *gocui.Gui, v1, v2 views.View) error {
	maxX, maxY := g.Size()
	err := v1.Delete(g)
	if err != nil {
		return err
	}

	err = v2.Set(g, 0, 6, maxX-1, maxY)
	if err != nil {
		return err
	}

	_, err = g.SetCurrentView(v2.Name())
	return err
}

func newController(app *app.App) *controller {
	m := models.New(app)
	return &controller{
		logger: app.Logger.With(logging.String("logger", "controller")),
		models: m,
		views:  views.New(app.Config.Views, m),
	}
}
