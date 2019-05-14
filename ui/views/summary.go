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
	SUMMARY_LEFT  = "summary_left"
	SUMMARY_RIGHT = "summary_right"
)

type Summary struct {
	left            *gocui.View
	right           *gocui.View
	info            *models.Info
	channelsBalance *models.ChannelsBalance
	walletBalance   *models.WalletBalance
	channels        *models.Channels
}

func (s *Summary) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	s.left, err = g.SetView(SUMMARY_LEFT, x0, y0, x1/2, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	s.left.Frame = false
	s.left.Wrap = true

	s.right, err = g.SetView(SUMMARY_RIGHT, x1/2, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	s.right.Frame = false
	s.right.Wrap = true
	s.display()
	return nil
}

func (s *Summary) display() {
	s.left.Clear()
	p := message.NewPrinter(language.English)
	green := color.Green()
	yellow := color.Yellow()
	cyan := color.Cyan()
	red := color.Red()
	fmt.Fprintln(s.left, green("[ Channels ]"))
	fmt.Fprintln(s.left, p.Sprintf("%s %d (%s|%s)",
		cyan("balance:"),
		s.channelsBalance.Balance+s.channelsBalance.PendingOpenBalance,
		green(p.Sprintf("%d", s.channelsBalance.Balance)),
		yellow(p.Sprintf("%d", s.channelsBalance.PendingOpenBalance)),
	))
	fmt.Fprintln(s.left, fmt.Sprintf("%s %d %s %d %s %d %s",
		cyan("state  :"),
		s.info.NumActiveChannels, green("active"),
		s.info.NumPendingChannels, yellow("pending"),
		s.info.NumInactiveChannels, red("inactive"),
	))
	fmt.Fprintln(s.left, fmt.Sprintf("%s %s",
		cyan("gauge  :"),
		gaugeTotal(s.channelsBalance.Balance, s.channels.List()),
	))

	s.right.Clear()
	fmt.Fprintln(s.right, green("[ Wallet ]"))
	fmt.Fprintln(s.right, p.Sprintf("%s %d (%s|%s)",
		cyan("balance:"),
		s.walletBalance.TotalBalance,
		green(p.Sprintf("%d", s.walletBalance.ConfirmedBalance)),
		yellow(p.Sprintf("%d", s.walletBalance.UnconfirmedBalance)),
	))
}

func gaugeTotal(balance int64, channels []*netmodels.Channel) string {
	capacity := int64(0)
	for i := range channels {
		capacity += channels[i].Capacity
	}

	if capacity == 0 {
		return fmt.Sprintf("[%20s]  0%%", "")
	}

	index := int(balance * int64(20) / capacity)
	var buffer bytes.Buffer
	cyan := color.Cyan()
	for i := 0; i < 20; i++ {
		if i < index {
			buffer.WriteString(cyan("|"))
			continue
		}
		buffer.WriteString(" ")
	}
	return fmt.Sprintf("[%s] %2d%%", buffer.String(), balance*100/capacity)
}

func NewSummary(info *models.Info,
	channelsBalance *models.ChannelsBalance,
	walletBalance *models.WalletBalance,
	channels *models.Channels) *Summary {
	return &Summary{
		info:            info,
		channelsBalance: channelsBalance,
		walletBalance:   walletBalance,
		channels:        channels,
	}
}
