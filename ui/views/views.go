package views

import (
	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/ui/models"
	"github.com/jroimartin/gocui"
)

type view interface {
	Set(*gocui.Gui, int, int, int, int) error
	Wrap(*gocui.View) view
	CursorLeft() error
	CursorRight() error
	CursorUp() error
	CursorDown() error
	Name() string
}

type Views struct {
	Previous view
	Main     view

	Help     *Help
	Header   *Header
	Menu     *Menu
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
	case MENU:
		return v.Menu.Wrap(vi)
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

	current := g.CurrentView()
	if current != nil && current.Name() == v.Menu.Name() {
		err = v.Menu.Set(g, 0, 6, 10, maxY)
		if err != nil {
			return err
		}

		err = v.Main.Set(g, 11, 6, maxX-1, maxY)
		if err != nil {
			return err
		}
		return nil
	}

	err = v.Main.Set(g, 0, 6, maxX-1, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	return nil
}

func New(cfg config.Views, m *models.Models) *Views {
	main := NewChannels(cfg.Channels, m.Channels)
	return &Views{
		Header:   NewHeader(m.Info),
		Help:     NewHelp(),
		Menu:     NewMenu(),
		Summary:  NewSummary(m.Info, m.ChannelsBalance, m.WalletBalance, m.Channels),
		Channels: main,
		Channel:  NewChannel(m.CurrentChannel),
		Main:     main,
	}
}
