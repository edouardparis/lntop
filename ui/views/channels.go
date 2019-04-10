package views

import (
	"bytes"
	"fmt"

	"github.com/jroimartin/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	netmodels "github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
)

const (
	CHANNEL          = "channel"
	CHANNEL_HEADER   = "channel_header"
	CHANNEL_FOOTER   = "channel_footer"
	CHANNELS         = "channels"
	CHANNELS_COLUMNS = "channels_columns"
	CHANNELS_FOOTER  = "channels_footer"
)

type Channels struct {
	view     *gocui.View
	channels *models.Channels
}

func (c Channels) Name() string {
	return CHANNELS
}

func (c *Channels) Wrap(v *gocui.View) view {
	c.view = v
	return c
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
	columns.FgColor = gocui.ColorBlack
	displayChannelsColumns(columns)

	c.view, err = g.SetView(CHANNELS, x0-1, y0+1, x1+2, y1-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		_, err = g.SetCurrentView(CHANNELS)
		if err != nil {
			return err
		}
	}
	c.view.Frame = false
	c.view.Autoscroll = false
	c.view.SelBgColor = gocui.ColorCyan
	c.view.SelFgColor = gocui.ColorBlack
	c.view.Highlight = true

	c.display()

	footer, err := g.SetView(CHANNELS_FOOTER, x0-1, y1-2, x1+2, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	footer.Frame = false
	footer.BgColor = gocui.ColorCyan
	footer.FgColor = gocui.ColorBlack
	footer.Clear()
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s %s%s",
		color.BlackBg("F1"), "Help",
		color.BlackBg("Enter"), "Channel",
		color.BlackBg("F10"), "Quit",
	))
	return nil
}

func displayChannelsColumns(v *gocui.View) {
	v.Clear()
	fmt.Fprintln(v, fmt.Sprintf("%-13s %-20s %-21s %12s %12s %5s  %-15s %s",
		"STATUS",
		"ALIAS",
		"GAUGE",
		"LOCAL",
		"CAP",
		"HTLC",
		"Last Update",
		"ID",
	))
}

func (c *Channels) display() {
	p := message.NewPrinter(language.English)
	c.view.Clear()
	for _, item := range c.channels.List() {
		line := fmt.Sprintf("%s %-20s %s %s %s %5d  %15s %d %500s",
			status(item),
			alias(item),
			gauge(item),
			color.Cyan(p.Sprintf("%12d", item.LocalBalance)),
			p.Sprintf("%12d", item.Capacity),
			len(item.PendingHTLC),
			lastUpdate(item),
			item.ID,
			"",
		)
		fmt.Fprintln(c.view, line)
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

func status(c *netmodels.Channel) string {
	switch c.Status {
	case netmodels.ChannelActive:
		return color.Green(fmt.Sprintf("%-13s", "active"))
	case netmodels.ChannelInactive:
		return color.Red(fmt.Sprintf("%-13s", "inactive"))
	case netmodels.ChannelOpening:
		return color.Yellow(fmt.Sprintf("%-13s", "opening"))
	case netmodels.ChannelClosing:
		return color.Yellow(fmt.Sprintf("%-13s", "closing"))
	case netmodels.ChannelForceClosing:
		return color.Yellow(fmt.Sprintf("%-13s", "force closing"))
	case netmodels.ChannelWaitingClose:
		return color.Yellow(fmt.Sprintf("%-13s", "waiting close"))
	}
	return ""
}

func gauge(c *netmodels.Channel) string {
	index := int(c.LocalBalance * int64(15) / c.Capacity)
	var buffer bytes.Buffer
	for i := 0; i < 15; i++ {
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
	view    *gocui.View
	channel *models.Channel
}

func (c Channel) Name() string {
	return CHANNEL
}

func (c Channel) Empty() bool {
	return c.channel == nil
}

func (c *Channel) Wrap(v *gocui.View) view {
	c.view = v
	return c
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

	v, err := g.SetView(CHANNEL, x0-1, y0+1, x1+2, y1-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Frame = false
	c.view = v
	c.display()

	footer, err := g.SetView(CHANNEL_FOOTER, x0-1, y1-2, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	footer.Frame = false
	footer.BgColor = gocui.ColorCyan
	footer.FgColor = gocui.ColorBlack
	footer.Clear()
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s %s%s",
		color.BlackBg("F1"), "Help",
		color.BlackBg("Enter"), "Channels",
		color.BlackBg("F10"), "Quit",
	))
	return nil
}

func (c Channel) Delete(g *gocui.Gui) error {
	err := g.DeleteView(CHANNEL_HEADER)
	if err != nil {
		return err
	}

	err = g.DeleteView(CHANNEL)
	if err != nil {
		return err
	}

	return g.DeleteView(CHANNEL_FOOTER)
}

func (c *Channel) display() {
	p := message.NewPrinter(language.English)
	v := c.view
	v.Clear()
	channel := c.channel.Item
	fmt.Fprintln(v, color.Green(" [ Channel ]"))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		color.Cyan("         Status:"), status(channel)))
	fmt.Fprintln(v, fmt.Sprintf("%s %d",
		color.Cyan("             ID:"), channel.ID))
	fmt.Fprintln(v, p.Sprintf("%s %d",
		color.Cyan("       Capacity:"), channel.Capacity))
	fmt.Fprintln(v, p.Sprintf("%s %d",
		color.Cyan("  Local Balance:"), channel.LocalBalance))
	fmt.Fprintln(v, p.Sprintf("%s %d",
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
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan(" Total Capacity:"), channel.Node.TotalCapacity))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan(" Total Channels:"), channel.Node.NumChannels))
	}

	if channel.Policy1 != nil {
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, color.Green(" [ Forward Policy Node1 ]"))
		if channel.Policy1.Disabled {
			fmt.Fprintln(v, color.Red("disabled"))
		}
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan("    Time lock delta:"), channel.Policy1.TimeLockDelta))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan("           Min htlc:"), channel.Policy1.MinHtlc))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan("      Fee base msat:"), channel.Policy1.FeeBaseMsat))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan("Fee rate milli msat:"), channel.Policy1.FeeRateMilliMsat))
	}

	if channel.Policy2 != nil {
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, color.Green(" [ Forward Policy Node 2 ]"))
		if channel.Policy2.Disabled {
			fmt.Fprintln(v, color.Red("disabled"))
		}
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan("    Time lock delta:"), channel.Policy2.TimeLockDelta))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan("           Min htlc:"), channel.Policy2.MinHtlc))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan("      Fee base msat:"), channel.Policy2.FeeBaseMsat))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			color.Cyan("Fee rate milli msat:"), channel.Policy2.FeeRateMilliMsat))
	}
}

func NewChannel(channel *models.Channel) *Channel {
	return &Channel{channel: channel}
}
