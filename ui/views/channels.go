package views

import (
	"bytes"
	"fmt"

	"github.com/jroimartin/gocui"

	netmodels "github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
)

const (
	CHANNEL          = "channel"
	CHANNELS         = "channels"
	CHANNELS_COLUMNS = "channels_columns"
)

type Channels struct {
	channels *models.Channels
}

func (c Channels) Name() string {
	return CHANNELS
}

func (c *Channels) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	columns, err := g.SetView(CHANNELS_COLUMNS, x0-1, y0, x1+2, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	columns.Frame = false
	columns.BgColor = gocui.ColorGreen
	columns.FgColor = gocui.ColorBlack | gocui.AttrBold
	displayChannelsColumns(columns)

	v, err := g.SetView(CHANNELS, x0-1, y0+1, x1+2, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		_, err = g.SetCurrentView(CHANNELS)
		if err != nil {
			return err
		}
	}
	v.Frame = false
	v.Autoscroll = true
	v.SelBgColor = gocui.ColorCyan
	v.SelFgColor = gocui.ColorBlack
	v.Highlight = true

	c.display(v)
	return nil
}

func displayChannelsColumns(v *gocui.View) {
	fmt.Fprintln(v, fmt.Sprintf("%-9s %-26s %12s %12s %5s",
		"Status",
		"Gauge",
		"Local",
		"Capacity",
		"pHTLC",
	))
}

func (c *Channels) display(v *gocui.View) {
	v.Clear()
	for _, item := range c.channels.Items {
		line := fmt.Sprintf("%s %s %s %12d %5d %500s",
			active(item),
			gauge(item),
			color.Cyan(fmt.Sprintf("%12d", item.LocalBalance)),
			item.Capacity,
			len(item.PendingHTLC),
			"",
		)
		fmt.Fprintln(v, line)
	}
}

func active(c *netmodels.Channel) string {
	if c.Active {
		return color.Green(fmt.Sprintf("%-9s", "active"))
	}
	return color.Red(fmt.Sprintf("%-9s", "inactive"))
}

func gauge(c *netmodels.Channel) string {
	index := int(c.LocalBalance * int64(20) / c.Capacity)
	var buffer bytes.Buffer
	for i := 0; i < 20; i++ {
		if i < index {
			buffer.WriteString(color.Cyan("|"))
			continue
		}
		buffer.WriteString(" ")
	}
	return fmt.Sprintf("[%s] %2d%%", buffer.String(), c.LocalBalance*100/c.Capacity)
}

func NewChannels(channels *models.Channels) *Channels {
	return &Channels{channels: channels}
}

type Channel struct {
	channel *models.Channel
}

func (c Channel) Name() string {
	return CHANNEL
}

func (c Channel) Empty() bool {
	return c.channel == nil
}

func (c *Channel) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(CHANNEL, x0-1, y0, x1+2, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	c.display(v)
	return nil
}

func (c *Channel) display(v *gocui.View) {
	v.Clear()
	fmt.Fprintln(v, fmt.Sprintf("%s %d",
		color.Cyan("ID:"), c.channel.Item.ID))
}

func NewChannel(channel *models.Channel) *Channel {
	return &Channel{channel: channel}
}
