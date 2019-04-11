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
	columns  *gocui.View
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

func (c *Channels) CursorDown() error {
	return cursorDown(c.view, 1)
}

func (c *Channels) CursorUp() error {
	return cursorUp(c.view, 1)
}

func (c *Channels) CursorRight() error {
	err := cursorRight(c.columns, 2)
	if err != nil {
		return err
	}

	return cursorRight(c.view, 2)
}

func (c *Channels) CursorLeft() error {
	err := cursorLeft(c.columns, 2)
	if err != nil {
		return err
	}

	return cursorLeft(c.view, 2)
}

func (c *Channels) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	c.columns, err = g.SetView(CHANNELS_COLUMNS, x0-1, y0, x1+2, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	c.columns.Frame = false
	c.columns.BgColor = gocui.ColorGreen
	c.columns.FgColor = gocui.ColorBlack
	displayChannelsColumns(c.columns)

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
	fmt.Fprintln(v, fmt.Sprintf("%-13s %-25s %-21s %12s %12s %5s %-10s %-6s %-15s %s %-19s",
		"STATUS",
		"ALIAS",
		"GAUGE",
		"LOCAL",
		"CAP",
		"HTLC",
		"UNSETTLED",
		"CFEE",
		"Last Update",
		"PRIVATE",
		"ID",
	))
}

func (c *Channels) display() {
	p := message.NewPrinter(language.English)
	c.view.Clear()
	for _, item := range c.channels.List() {
		line := fmt.Sprintf("%s %-25s %s %s %s %5d %s %s %s %s %19s %500s",
			status(item),
			alias(item),
			gauge(item),
			color.Cyan(p.Sprintf("%12d", item.LocalBalance)),
			p.Sprintf("%12d", item.Capacity),
			len(item.PendingHTLC),
			color.Yellow(p.Sprintf("%10d", item.UnsettledBalance)),
			p.Sprintf("%6d", item.CommitFee),
			lastUpdate(item),
			channelPrivate(item),
			channelID(item),
			"",
		)
		fmt.Fprintln(c.view, line)
	}
}

func channelPrivate(c *netmodels.Channel) string {
	if c.Private {
		return color.Red("private")
	}

	return color.Green("public ")
}

func channelID(c *netmodels.Channel) string {
	if c.ID == 0 {
		return ""
	}

	return fmt.Sprintf("%d", c.ID)
}

func alias(c *netmodels.Channel) string {
	if c.Node == nil || c.Node.Alias == "" {
		return c.RemotePubKey[:19]
	}

	return c.Node.Alias
}

func lastUpdate(c *netmodels.Channel) string {
	if c.LastUpdate != nil {
		return color.Cyan(
			fmt.Sprintf("%15s", c.LastUpdate.Format("15:04:05 Jan _2")),
		)
	}

	return fmt.Sprintf("%15s", "")
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

func (c *Channel) CursorDown() error {
	return cursorDown(c.view, 1)
}

func (c *Channel) CursorUp() error {
	return cursorUp(c.view, 1)
}

func (c *Channel) CursorRight() error {
	return cursorRight(c.view, 1)
}

func (c *Channel) CursorLeft() error {
	return cursorLeft(c.view, 1)
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
