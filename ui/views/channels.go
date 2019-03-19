package views

import (
	"context"
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/network"
	"github.com/edouardparis/lntop/network/models"
)

const CHANNELS = "channels"

type Channels struct {
	*gocui.View
	items   []*models.Channel
	network *network.Network
}

func (c *Channels) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	c.View, err = g.SetView(CHANNELS, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	c.View.Frame = false
	err = c.Update(context.Background())
	if err != nil {
		return err
	}

	c.Display()
	return nil
}

func (c *Channels) Update(ctx context.Context) error {
	channels, err := c.network.ListChannels(ctx)
	if err != nil {
		return err
	}

	c.items = channels
	return nil
}

func (c *Channels) Display() {
	for i := range c.items {
		fmt.Fprintln(c.View, fmt.Sprintf("%d", c.items[i].ID))
	}
}

func NewChannels(network *network.Network) *Channels {
	return &Channels{network: network}
}
