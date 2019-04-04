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
	CHANNEL_HEADER   = "channel_header"
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
	v.Clear()
	fmt.Fprintln(v, fmt.Sprintf("%-9s %-20s %-26s %12s %12s %5s  %-15s %s",
		"Status",
		"Alias",
		"Gauge",
		"Local",
		"Capacity",
		"pHTLC",
		"Last Update",
		"ID",
	))
}

func (c *Channels) display(v *gocui.View) {
	v.Clear()
	for _, item := range c.channels.List() {
		line := fmt.Sprintf("%s %-20s %s %s %12d %5d  %15s %d %500s",
			active(item),
			alias(item),
			gauge(item),
			color.Cyan(fmt.Sprintf("%12d", item.LocalBalance)),
			item.Capacity,
			len(item.PendingHTLC),
			lastUpdate(item),
			item.ID,
			"",
		)
		fmt.Fprintln(v, line)
	}
}

func alias(c *netmodels.Channel) string {
	if c.Node == nil || c.Node.Alias == "" {
		return c.RemotePubKey[:19]
	}

	return c.Node.Alias
}

func lastUpdate(c *netmodels.Channel) string {
	if c.LastUpdate != nil {
		return c.LastUpdate.Format("15:04:05 Jan _2")
	}

	return ""
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
	header, err := g.SetView(CHANNEL_HEADER, x0-1, y0, x1+2, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	header.Frame = false
	header.BgColor = gocui.ColorGreen
	header.FgColor = gocui.ColorBlack | gocui.AttrBold
	header.Clear()
	fmt.Fprintln(header, "Channel")

	v, err := g.SetView(CHANNEL, x0-1, y0+1, x1+2, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Frame = false
	c.display(v)
	return nil
}

func (c Channel) Delete(g *gocui.Gui) error {
	err := g.DeleteView(CHANNEL_HEADER)
	if err != nil {
		return err
	}
	return g.DeleteView(CHANNEL)
}

func (c *Channel) display(v *gocui.View) {
	v.Clear()
	channel := c.channel.Item
	fmt.Fprintln(v, color.Green(" [ Channel ]"))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		color.Cyan("         Status:"), active(channel)))
	fmt.Fprintln(v, fmt.Sprintf("%s %d",
		color.Cyan("             ID:"), channel.ID))
	fmt.Fprintln(v, fmt.Sprintf("%s %d",
		color.Cyan("       Capacity:"), channel.Capacity))
	fmt.Fprintln(v, fmt.Sprintf("%s %d",
		color.Cyan("  Local Balance:"), channel.LocalBalance))
	fmt.Fprintln(v, fmt.Sprintf("%s %d",
		color.Cyan(" Remote Balance:"), channel.RemoteBalance))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		color.Cyan("  Channel Point:"), channel.ChannelPoint))
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, color.Green(" [ Node ]"))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		color.Cyan("          Alias:"), alias(channel)))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		color.Cyan("         PubKey:"), channel.RemotePubKey))

	if channel.Node != nil {
		fmt.Fprintln(v, fmt.Sprintf("%s %d",
			color.Cyan(" Total Capacity:"), channel.Node.TotalCapacity))
		fmt.Fprintln(v, fmt.Sprintf("%s %d",
			color.Cyan(" Total Channels:"), channel.Node.NumChannels))
	}
}

func NewChannel(channel *models.Channel) *Channel {
	return &Channel{channel: channel}
}
