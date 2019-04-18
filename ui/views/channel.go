package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
)

const (
	CHANNEL        = "channel"
	CHANNEL_HEADER = "channel_header"
	CHANNEL_FOOTER = "channel_footer"
)

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

	v, err := g.SetView(CHANNEL, x0-1, y0+1, x1+2, y1-1)
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
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s %s%s %s%s",
		color.BlackBg("F1"), "Help",
		color.BlackBg("F2"), "Menu",
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
		color.Cyan("         PubKey:"), channel.RemotePubKey))

	if channel.Node != nil {
		fmt.Fprintln(v, fmt.Sprintf("%s %s",
			color.Cyan("          Alias:"), channel.Node.Alias))
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
