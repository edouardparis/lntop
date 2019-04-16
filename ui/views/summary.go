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
	fmt.Fprintln(s.left, color.Green("[ Channels ]"))
	fmt.Fprintln(s.left, p.Sprintf("%s %d (%s|%s)",
		color.Cyan("balance:"),
		s.channelsBalance.Balance+s.channelsBalance.PendingOpenBalance,
		color.Green(p.Sprintf("%d", s.channelsBalance.Balance)),
		color.Yellow(p.Sprintf("%d", s.channelsBalance.PendingOpenBalance)),
	))
	fmt.Fprintln(s.left, fmt.Sprintf("%s %d %s %d %s %d %s",
		color.Cyan("state  :"),
		s.info.NumActiveChannels, color.Green("active"),
		s.info.NumPendingChannels, color.Yellow("pending"),
		s.info.NumInactiveChannels, color.Red("inactive"),
	))
	fmt.Fprintln(s.left, fmt.Sprintf("%s %s",
		color.Cyan("gauge  :"),
		gaugeTotal(s.channelsBalance.Balance, s.channels.List()),
	))

	s.right.Clear()
	fmt.Fprintln(s.right, color.Green("[ Wallet ]"))
	fmt.Fprintln(s.right, p.Sprintf("%s %d (%s|%s)",
		color.Cyan("balance:"),
		s.walletBalance.TotalBalance,
		color.Green(p.Sprintf("%d", s.walletBalance.ConfirmedBalance)),
		color.Yellow(p.Sprintf("%d", s.walletBalance.UnconfirmedBalance)),
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
	for i := 0; i < 20; i++ {
		if i < index {
			buffer.WriteString(color.Cyan("|"))
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
