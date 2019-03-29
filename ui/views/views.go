package views

import (
	"github.com/edouardparis/lntop/ui/models"
	"github.com/jroimartin/gocui"
)

type Views struct {
	Header   *Header
	Summary  *Summary
	Channels *Channels
	Footer   *Footer
}

func (v *Views) Layout(g *gocui.Gui, maxX, maxY int) error {
	err := v.Header.Set(g, 0, -1, maxX, 1)
	if err != nil {
		return err
	}

	err = v.Summary.Set(g, 0, 1, maxX, 6)
	if err != nil {
		return err
	}

	err = v.Channels.Set(g, 0, 6, maxX-1, maxY-1)
	if err != nil {
		return err
	}

	return v.Footer.Set(g, 0, maxY-2, maxX, maxY)
}

func New(m *models.Models) *Views {
	return &Views{
		Header:   NewHeader(m.Info),
		Footer:   NewFooter(),
		Summary:  NewSummary(),
		Channels: NewChannels(m.Channels),
	}
}
