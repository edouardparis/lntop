package views

import (
	"github.com/edouardparis/lntop/ui/models"
	"github.com/jroimartin/gocui"
)

type view interface {
	Set(*gocui.Gui, int, int, int, int) error
	Wrap(*gocui.View) view
	Name() string
}

type Views struct {
	Previous view
	Help     *Help
	Header   *Header
	Summary  *Summary
	Channels *Channels
	Channel  *Channel
}

func (v Views) Get(vi *gocui.View) view {
	if vi == nil {
		return nil
	}
	switch vi.Name() {
	case CHANNELS:
		return v.Channels.Wrap(vi)
	case HELP:
		return v.Help.Wrap(vi)
	case CHANNEL:
		return v.Channel.Wrap(vi)
	default:
		return nil
	}
}

func (v *Views) SetPrevious(p view) {
	v.Previous = p
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

	return v.Channels.Set(g, 0, 6, maxX-1, maxY)
}

func New(m *models.Models) *Views {
	return &Views{
		Header:   NewHeader(m.Info),
		Help:     NewHelp(),
		Summary:  NewSummary(m.Info, m.ChannelsBalance, m.WalletBalance, m.Channels),
		Channels: NewChannels(m.Channels),
		Channel:  NewChannel(m.CurrentChannel),
	}
}
