package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	netmodels "github.com/edouardparis/lntop/network/models"
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
	fmt.Fprintf(footer, "%s%s %s%s %s%s %s%s\n",
		blackBg("F2"), "Menu",
		blackBg("Enter"), "Channels",
		blackBg("C"), "Get disabled",
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

func printPolicy(v *gocui.View, p *message.Printer, policy *netmodels.RoutingPolicy, outgoing bool) {
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
	fmt.Fprintf(v, "%s %s\n",
		cyan("     Min htlc (msat):"), formatAmount(policy.MinHtlc))
	fmt.Fprintf(v, "%s %s\n",
		cyan("      Max htlc (sat):"), formatAmount(int64(policy.MaxHtlc/1000)))
	fmt.Fprintf(v, "%s %s\n",
		cyan("       Fee base msat:"), formatAmount(policy.FeeBaseMsat))
	fmt.Fprintf(v, "%s %d\n",
		cyan(" Fee rate milli msat:"), policy.FeeRateMilliMsat)
}

func formatAmount(amt int64) string {
	btc := amt / 1e8
	ms := amt % 1e8 / 1e6
	ts := amt % 1e6 / 1e3
	s := amt % 1e3
	if btc > 0 {
		return fmt.Sprintf("%d.%02d,%03d,%03d", btc, ms, ts, s)
	}
	if ms > 0 {
		return fmt.Sprintf("%d,%03d,%03d", ms, ts, s)
	}
	if ts > 0 {
		return fmt.Sprintf("%d,%03d", ts, s)
	}
	if s >= 0 {
		return fmt.Sprintf("%d", s)
	}
	return fmt.Sprintf("error: %d", amt)
}

func formatDisabledCount(cnt int, total uint32) string {
	perc := uint32(cnt) * 100 / total
	disabledStr := ""
	if perc >= 25 && perc < 50 {
		disabledStr = color.Yellow(color.Bold)(fmt.Sprintf("%4d", cnt))
	} else if perc >= 50 {
		disabledStr = color.Red(color.Bold)(fmt.Sprintf("%4d", cnt))
	} else {
		disabledStr = fmt.Sprintf("%4d", cnt)
	}
	return fmt.Sprintf("%s / %d (%d%%)", disabledStr, total, perc)
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
		cyan("             Status:"), status(channel))
	if channel.Status == netmodels.ChannelForceClosing {
		fmt.Fprintf(v, "%s %d blocks\n",
			cyan("         Matured in:"), channel.BlocksTilMaturity)
	}
	fmt.Fprintf(v, "%s %d (%s)\n",
		cyan("                 ID:"), channel.ID, ToScid(channel.ID))
	fmt.Fprintf(v, "%s %s\n",
		cyan("           Capacity:"), formatAmount(channel.Capacity))
	fmt.Fprintf(v, "%s %s\n",
		cyan("      Local Balance:"), formatAmount(channel.LocalBalance))
	fmt.Fprintf(v, "%s %s\n",
		cyan("     Remote Balance:"), formatAmount(channel.RemoteBalance))
	fmt.Fprintf(v, "%s %s\n",
		cyan("      Channel Point:"), channel.ChannelPoint)
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
		fmt.Fprintf(v, "%s %s\n",
			cyan(" Total Capacity:"), formatAmount(channel.Node.TotalCapacity))
		fmt.Fprintf(v, "%s %d\n",
			cyan(" Total Channels:"), channel.Node.NumChannels)

		if c.channels.CurrentNode != nil && c.channels.CurrentNode.PubKey == channel.RemotePubKey {
			disabledOut := 0
			disabledIn := 0
			for _, ch := range c.channels.CurrentNode.Channels {
				if ch.LocalPolicy != nil && ch.LocalPolicy.Disabled {
					disabledOut++
				}
				if ch.RemotePolicy != nil && ch.RemotePolicy.Disabled {
					disabledIn++
				}
			}
			fmt.Fprintf(v, "\n %s %s\n", cyan("Disabled from node:"), formatDisabledCount(disabledOut, channel.Node.NumChannels))
			fmt.Fprintf(v, " %s %s\n", cyan("Disabled to node:  "), formatDisabledCount(disabledIn, channel.Node.NumChannels))
		}
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
			fmt.Fprintf(v, "%s %s\n",
				cyan("     Amount:"), formatAmount(htlc.Amount))
			fmt.Fprintf(v, "%s %d\n",
				cyan(" Expiration:"), htlc.ExpirationHeight)
			fmt.Fprintln(v)
		}
	}

}

func NewChannel(channels *models.Channels) *Channel {
	return &Channel{channels: channels}
}
