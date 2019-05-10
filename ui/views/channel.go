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
	view     *gocui.View
	channels *models.Channels
}

func (c Channel) Name() string {
	return CHANNEL
}

func (c Channel) Empty() bool {
	return c.channels == nil
}

func (c *Channel) Wrap(v *gocui.View) View {
	c.view = v
	return c
}

func (c Channel) Origin() (int, int) {
	return c.view.Origin()
}

func (c Channel) Cursor() (int, int) {
	return c.view.Cursor()
}

func (c Channel) Speed() (int, int, int, int) {
	return 1, 1, 1, 1
}

func (c *Channel) SetCursor(x, y int) error {
	return c.view.SetCursor(x, y)
}

func (c *Channel) SetOrigin(x, y int) error {
	return c.view.SetOrigin(x, y)
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
	blackBg := color.Black(color.Background)
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s %s%s %s%s",
		blackBg("F1"), "Help",
		blackBg("F2"), "Menu",
		blackBg("Enter"), "Channels",
		blackBg("F10"), "Quit",
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
	channel := c.channels.Current()
	green := color.Green()
	cyan := color.Cyan()
	red := color.Red()
	fmt.Fprintln(v, green(" [ Channel ]"))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		cyan("         Status:"), status(channel)))
	fmt.Fprintln(v, fmt.Sprintf("%s %d",
		cyan("             ID:"), channel.ID))
	fmt.Fprintln(v, p.Sprintf("%s %d",
		cyan("       Capacity:"), channel.Capacity))
	fmt.Fprintln(v, p.Sprintf("%s %d",
		cyan("  Local Balance:"), channel.LocalBalance))
	fmt.Fprintln(v, p.Sprintf("%s %d",
		cyan(" Remote Balance:"), channel.RemoteBalance))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		cyan("  Channel Point:"), channel.ChannelPoint))
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, green(" [ Node ]"))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		cyan("         PubKey:"), channel.RemotePubKey))

	if channel.Node != nil {
		fmt.Fprintln(v, fmt.Sprintf("%s %s",
			cyan("          Alias:"), channel.Node.Alias))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan(" Total Capacity:"), channel.Node.TotalCapacity))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan(" Total Channels:"), channel.Node.NumChannels))
	}

	if channel.Policy1 != nil {
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, green(" [ Forward Policy Node1 ]"))
		if channel.Policy1.Disabled {
			fmt.Fprintln(v, red("disabled"))
		}
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan("    Time lock delta:"), channel.Policy1.TimeLockDelta))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan("           Min htlc:"), channel.Policy1.MinHtlc))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan("      Fee base msat:"), channel.Policy1.FeeBaseMsat))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan("Fee rate milli msat:"), channel.Policy1.FeeRateMilliMsat))
	}

	if channel.Policy2 != nil {
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, green(" [ Forward Policy Node 2 ]"))
		if channel.Policy2.Disabled {
			fmt.Fprintln(v, red("disabled"))
		}
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan("    Time lock delta:"), channel.Policy2.TimeLockDelta))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan("           Min htlc:"), channel.Policy2.MinHtlc))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan("      Fee base msat:"), channel.Policy2.FeeBaseMsat))
		fmt.Fprintln(v, p.Sprintf("%s %d",
			cyan("Fee rate milli msat:"), channel.Policy2.FeeRateMilliMsat))
	}
}

func NewChannel(channels *models.Channels) *Channel {
	return &Channel{channels: channels}
}
