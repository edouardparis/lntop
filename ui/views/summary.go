package views

import (
	"fmt"

	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
	"github.com/jroimartin/gocui"
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

	s.right, err = g.SetView(SUMMARY_RIGHT, x1/2, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	s.right.Frame = false
	s.display()
	return nil
}

func (s *Summary) display() {
	s.left.Clear()
	fmt.Fprintln(s.left, color.Green("[ Channels ]"))
	fmt.Fprintln(s.left, fmt.Sprintf("%s %d (%s|%s)",
		color.Cyan("balance:"),
		s.channelsBalance.Balance+s.channelsBalance.PendingOpenBalance,
		color.Green(fmt.Sprintf("%d", s.channelsBalance.Balance)),
		color.Yellow(fmt.Sprintf("%d", s.channelsBalance.PendingOpenBalance)),
	))
	fmt.Fprintln(s.left, fmt.Sprintf("%s %d %s %d %s %d %s",
		color.Cyan("state  :"),
		s.info.NumActiveChannels, color.Green("active"),
		s.info.NumPendingChannels, color.Yellow("pending"),
		s.info.NumInactiveChannels, color.Red("inactive"),
	))

	s.right.Clear()
	fmt.Fprintln(s.right, color.Green("[ Wallet ]"))
	fmt.Fprintln(s.right, fmt.Sprintf("%s %d (%s|%s)",
		color.Cyan("balance:"),
		s.walletBalance.TotalBalance,
		color.Green(fmt.Sprintf("%d", s.walletBalance.ConfirmedBalance)),
		color.Yellow(fmt.Sprintf("%d", s.walletBalance.UnconfirmedBalance)),
	))
}

func NewSummary(info *models.Info, channelsBalance *models.ChannelsBalance, walletBalance *models.WalletBalance) *Summary {
	return &Summary{
		info:            info,
		channelsBalance: channelsBalance,
		walletBalance:   walletBalance,
	}
}
