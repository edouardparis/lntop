package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"github.com/pkg/errors"

	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/ui/cursor"
	"github.com/edouardparis/lntop/ui/models"
)

type View interface {
	Set(*gocui.Gui, int, int, int, int) error
	Delete(*gocui.Gui) error
	Wrap(*gocui.View) View
	Name() string
	cursor.View
}

type Views struct {
	Main View

	Header       *Header
	Menu         *Menu
	Summary      *Summary
	Channels     *Channels
	Channel      *Channel
	Transactions *Transactions
	Transaction  *Transaction
	Routing      *Routing
}

func (v Views) Get(vi *gocui.View) View {
	if vi == nil {
		return nil
	}
	switch vi.Name() {
	case CHANNELS:
		return v.Channels.Wrap(vi)
	case MENU:
		return v.Menu.Wrap(vi)
	case CHANNEL:
		return v.Channel.Wrap(vi)
	case TRANSACTIONS:
		return v.Transactions.Wrap(vi)
	case TRANSACTION:
		return v.Transaction.Wrap(vi)
	case ROUTING:
		return v.Routing.Wrap(vi)
	default:
		return nil
	}
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
	if current != nil {
		switch current.Name() {
		case v.Menu.Name():
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
	}

	err = v.Main.Set(g, 0, 6, maxX-1, maxY)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	_, err = g.SetCurrentView(v.Main.Name())
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func New(cfg config.Views, m *models.Models) *Views {
	main := NewChannels(cfg.Channels, m.Channels)
	return &Views{
		Header:       NewHeader(m.Info),
		Menu:         NewMenu(),
		Summary:      NewSummary(m.Info, m.ChannelsBalance, m.WalletBalance, m.Channels),
		Channels:     main,
		Channel:      NewChannel(m.Channels),
		Transactions: NewTransactions(cfg.Transactions, m.Transactions),
		Transaction:  NewTransaction(m.Transactions),
		Routing:      NewRouting(cfg.Routing, m.RoutingLog, m.Channels),
		Main:         main,
	}
}

func ToScid(id uint64) string {
	blocknum := id >> 40
	txnum := (id >> 16) & 0x00FFFFFF
	outnum := id & 0xFFFF

	return fmt.Sprintf("%dx%dx%d", blocknum, txnum, outnum)
}
