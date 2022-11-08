package views

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/pkg/errors"

	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/ui/color"
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
	FwdingHist   *FwdingHist
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
	case FWDINGHIST:
		return v.FwdingHist.Wrap(vi)
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
		if current.Name() == v.Menu.Name() {
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
		FwdingHist:   NewFwdingHist(cfg.FwdingHist, m.FwdingHist),
		Main:         main,
	}
}

func ToScid(id uint64) string {
	blocknum := id >> 40
	txnum := (id >> 16) & 0x00FFFFFF
	outnum := id & 0xFFFF

	return fmt.Sprintf("%dx%dx%d", blocknum, txnum, outnum)
}

func FormatAge(age uint32) string {
	if age < 6 {
		return fmt.Sprintf("%02dm", age*10)
	} else if age < 144 {
		return fmt.Sprintf("%02dh", age/6)
	} else if age < 4383 {
		return fmt.Sprintf("%02dd%02dh", age/144, (age%144)/6)
	} else if age < 52596 {
		return fmt.Sprintf("%02dm%02dd%02dh", age/4383, (age%4383)/144, (age%144)/6)
	}
	return fmt.Sprintf("%02dy%02dm%02dd", age/52596, (age%52596)/4383, (age%4383)/144)
}

func interp(a, b [3]float64, r float64) (result [3]float64) {
	result[0] = a[0] + (b[0]-a[0])*r
	result[1] = a[1] + (b[1]-a[1])*r
	result[2] = a[2] + (b[2]-a[2])*r
	return
}

func ColorizeAge(age uint32, text string, opts ...color.Option) string {
	ageColors := [][3]float64{
		{120, 0.9, 0.9},
		{60, 0.9, 0.6},
		{22, 1, 0.5},
	}
	cur := [3]float64{}
	if age < 26298 {
		cur = interp(ageColors[0], ageColors[1], float64(age)/26298)
	} else if age < 52596 {
		cur = interp(ageColors[1], ageColors[2], float64(age-26298)/26298)
	} else {
		cur = ageColors[2]
	}
	return color.HSL256(cur[0]/360, cur[1], cur[2], opts...)(text)
}

func cursorCompat(v *gocui.View, x, y int) error {
	maxX, maxY := v.Size()
	if x < 0 || x >= maxX || y < 0 || y >= maxY {
		return gocui.ErrInvalidPoint
	}
	return nil
}
