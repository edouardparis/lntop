package views

import (
	"bytes"
	"fmt"

	"github.com/jroimartin/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/edouardparis/lntop/config"
	netmodels "github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
)

const (
	CHANNELS         = "channels"
	CHANNELS_COLUMNS = "channels_columns"
	CHANNELS_FOOTER  = "channels_footer"
)

type Channels struct {
	index    int
	cfg      *config.View
	columns  *gocui.View
	view     *gocui.View
	channels *models.Channels
}

func (c Channels) Index() int {
	return c.index
}

func (c Channels) Name() string {
	return CHANNELS
}

func (c *Channels) Wrap(v *gocui.View) view {
	c.view = v
	return c
}

func (c *Channels) CursorDown() error {
	if c.channels.Len() <= c.index+1 {
		return nil
	}
	c.index++
	return cursorDown(c.view, 1)
}

func (c *Channels) CursorUp() error {
	if c.index > 0 {
		c.index--
	}
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

func NewChannels(cfg *config.View, channels *models.Channels) *Channels {
	return &Channels{cfg: cfg, channels: channels}
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
		return c.RemotePubKey[:24]
	} else if len(c.Node.Alias) > 25 {
		return c.Node.Alias[:24]
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
