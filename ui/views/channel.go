package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	netModels "github.com/edouardparis/lntop/network/models"
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

func (c Channel) Limits() (pageSize int, fullSize int) {
	_, pageSize = c.view.Size()
	fullSize = len(c.view.BufferLines()) - 1
	return
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
	fmt.Fprintf(footer, "%s%s %s%s %s%s\n",
		blackBg("F2"), "Menu",
		blackBg("Enter"), "Channels",
		blackBg("F10"), "Quit",
	)
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

func printPolicy(v *gocui.View, p *message.Printer, policy *netModels.RoutingPolicy, outgoing bool) {
	green := color.Green()
	cyan := color.Cyan()
	red := color.Red()
	fmt.Fprintln(v, "")
	direction := "Outgoing"
	if !outgoing {
		direction = "Incoming"
	}
	fmt.Fprintf(v, green(" [ %s Policy ]\n"), direction)
	if policy.Disabled {
		fmt.Fprintln(v, red("disabled"))
	}
	fmt.Fprintf(v, "%s %d\n",
		cyan("     Time lock delta:"), policy.TimeLockDelta)
	fmt.Fprintf(v, "%s %d\n",
		cyan("            Min htlc:"), policy.MinHtlc)
	fmt.Fprintf(v, "%s %d\n",
		cyan("       Fee base msat:"), policy.FeeBaseMsat)
	fmt.Fprintf(v, "%s %d\n",
		cyan(" Fee rate milli msat:"), policy.FeeRateMilliMsat)
}

func (c *Channel) display() {
	p := message.NewPrinter(language.English)
	v := c.view
	v.Clear()
	channel := c.channels.Current()
	green := color.Green()
	cyan := color.Cyan()

	fmt.Fprintln(v, green(" [ Channel ]"))
	fmt.Fprintf(v, "%s %s\n",
		cyan("         Status:"), status(channel))
	fmt.Fprintf(v, "%s %d (%s)\n",
		cyan("             ID:"), channel.ID, ToScid(channel.ID))
	fmt.Fprintf(v, "%s %d\n",
		cyan("       Capacity:"), channel.Capacity)
	fmt.Fprintf(v, "%s %d\n",
		cyan("  Local Balance:"), channel.LocalBalance)
	fmt.Fprintf(v, "%s %d\n",
		cyan(" Remote Balance:"), channel.RemoteBalance)
	fmt.Fprintf(v, "%s %s\n",
		cyan("  Channel Point:"), channel.ChannelPoint)

	fmt.Fprintln(v, "")

	fmt.Fprintln(v, green(" [ Node ]"))
	fmt.Fprintf(v, "%s %s\n",
		cyan("         PubKey:"), channel.RemotePubKey)

	if channel.Node != nil {
		alias, forced := channel.ShortAlias()
		if forced {
			alias = cyan(alias)
		}
		fmt.Fprintf(v, "%s %s\n",
			cyan("          Alias:"), alias)
		fmt.Fprintf(v, "%s %d\n",
			cyan(" Total Capacity:"), channel.Node.TotalCapacity)
		fmt.Fprintf(v, "%s %d\n",
			cyan(" Total Channels:"), channel.Node.NumChannels)
	}
	if channel.LocalPolicy != nil {
		printPolicy(v, p, channel.LocalPolicy, true)
	}
	if channel.RemotePolicy != nil {
		printPolicy(v, p, channel.RemotePolicy, false)
	}
	if len(channel.PendingHTLC) > 0 {
		fmt.Fprintln(v)
		fmt.Fprintln(v, green(" [ Pending HTLCs ]"))
		for _, htlc := range channel.PendingHTLC {
			fmt.Fprintf(v, "%s %t\n",
				cyan("   Incoming:"), htlc.Incoming)
			fmt.Fprintf(v, "%s %d\n",
				cyan("     Amount:"), htlc.Amount)
			fmt.Fprintf(v, "%s %d\n",
				cyan(" Expiration:"), htlc.ExpirationHeight)
			fmt.Fprintln(v)
		}
	}

}

func NewChannel(channels *models.Channels) *Channel {
	return &Channel{channels: channels}
}
