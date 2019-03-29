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
	left  *gocui.View
	right *gocui.View
	info  *models.Info
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

func (s *Summary) display() {
	s.left.Clear()
	fmt.Fprintln(s.left, color.Green("[ Channels ]"))
	fmt.Fprintln(s.left, fmt.Sprintf("%d %s %d %s %d %s",
		s.info.NumPendingChannels, color.Yellow("pending"),
		s.info.NumActiveChannels, color.Green("active"),
		s.info.NumInactiveChannels, color.Red("inactive"),
	))

	s.right.Clear()
}

func NewSummary(info *models.Info) *Summary {
	return &Summary{info: info}
}
