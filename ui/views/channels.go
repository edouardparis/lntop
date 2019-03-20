package views

import (
	"context"
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/network"
	"github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
)

const (
	CHANNELS        = "channels"
	CHANNELS_HEADER = "header"
)

type Channels struct {
	*gocui.View
	items   []*models.Channel
	network *network.Network
}

func (c *Channels) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	headerView, err := g.SetView(CHANNELS_HEADER, x0, y0, x1, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	headerView.Frame = false
	headerView.BgColor = gocui.ColorGreen

	c.View, err = g.SetView(CHANNELS, x0, y0+1, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	c.View.Frame = false

	err = c.update(context.Background())
	if err != nil {
		return err
	}

	c.display()
	return nil
}

func (c *Channels) Refresh(g *gocui.Gui) error {
	var err error
	c.View, err = g.View(CHANNELS)
	if err != nil {
		return err
	}

	err = c.update(context.Background())
	if err != nil {
		return err
	}

	c.display()
	return nil
}

func (c *Channels) update(ctx context.Context) error {
	channels, err := c.network.ListChannels(ctx)
	if err != nil {
		return err
	}

	c.items = channels
	return nil
}

func (c *Channels) display() {
	for _, item := range c.items {
		line := fmt.Sprintf("%9s %d  %12d %12d  %s",
			active(item),
			item.ID,
			item.LocalBalance,
			item.Capacity,
			item.RemotePubKey,
		)
		fmt.Fprintln(c.View, line)
	}
}

func NewChannels(network *network.Network) *Channels {
	return &Channels{network: network}
}

func active(c *models.Channel) string {
	if c.Active {
		return color.Green("active  ")
	}
	return color.Red("inactive")
}
