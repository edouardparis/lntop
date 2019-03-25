package views

import (
	"fmt"

	"github.com/edouardparis/lntop/ui/color"
	"github.com/jroimartin/gocui"
)

const (
	SUMMARY_LEFT  = "summary_left"
	SUMMARY_RIGHT = "summary_right"
)

type Summary struct {
	left                *gocui.View
	right               *gocui.View
	NumPendingChannels  uint32
	NumActiveChannels   uint32
	NumInactiveChannels uint32
	NumPeers            uint32
	BlockHeight         uint32
	Synced              bool
}

func (s *Summary) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	s.left, err = g.SetView(SUMMARY_LEFT, x0, y0, x0+40, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	s.left.Frame = false

	s.right, err = g.SetView(SUMMARY_RIGHT, x0+40, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	s.right.Frame = false
	s.display()
	return nil
}

func (s *Summary) UpdateChannelsStats(numPendingChannels, numActiveChannels, numInactiveChannels uint32) {
	s.NumActiveChannels = numActiveChannels
	s.NumInactiveChannels = numInactiveChannels
	s.NumPendingChannels = numPendingChannels
}

func (s *Summary) display() {
	s.left.Clear()
	fmt.Fprintln(s.left, color.Green("[ Channels ]"))
	fmt.Fprintln(s.left, fmt.Sprintf("%d %s %d %s %d %s",
		s.NumPendingChannels, color.Yellow("pending"),
		s.NumActiveChannels, color.Green("active"),
		s.NumInactiveChannels, color.Red("inactive"),
	))

	s.right.Clear()
	fmt.Fprintln(s.right, color.Green("[ Network ]"))
	fmt.Fprintln(s.right, fmt.Sprintf("%s %4d", color.Cyan("Block Height: "), s.BlockHeight))
	fmt.Fprintln(s.right, fmt.Sprintf("%s %4d", color.Cyan("Synced:  "), s.NumActiveChannels))
}

func NewSummary() *Summary {
	return &Summary{}
}
