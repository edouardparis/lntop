package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	MENU        = "menu"
	MENU_HEADER = "menu_header"
)

var menu = []string{
	"CHANNEL",
	"TRANSAC",
}

type Menu struct {
	view *gocui.View
}

func (h Menu) Name() string {
	return MENU
}

func (h *Menu) Wrap(v *gocui.View) View {
	h.view = v
	return h
}

func (h Menu) Origin() (int, int) {
	return h.view.Origin()
}

func (h Menu) Cursor() (int, int) {
	return h.view.Cursor()
}

func (h Menu) Speed() (int, int, int, int) {
	return 1, 1, 1, 1
}

func (h *Menu) SetCursor(x, y int) error {
	return h.view.SetCursor(x, y)
}

func (h *Menu) SetOrigin(x, y int) error {
	return h.view.SetOrigin(x, y)
}

func (h Menu) Current() string {
	_, y := h.view.Cursor()
	if y < len(menu) {
		switch menu[y] {
		case "CHANNEL":
			return CHANNELS
		case "TRANSAC":
			return TRANSACTIONS
		}
	}
	return ""
}

func (c Menu) Delete(g *gocui.Gui) error {
	err := g.DeleteView(MENU_HEADER)
	if err != nil {
		return err
	}

	return g.DeleteView(MENU)
}

func (h Menu) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	header, err := g.SetView(MENU_HEADER, x0-1, y0, x1, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	header.Frame = false
	header.BgColor = gocui.ColorGreen
	header.FgColor = gocui.ColorBlack

	header.Clear()
	fmt.Fprintln(header, " MENU")

	h.view, err = g.SetView(MENU, x0-1, y0+1, x1, y1-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}

	h.view.Frame = false
	h.view.Highlight = true
	h.view.SelBgColor = gocui.ColorCyan
	h.view.SelFgColor = gocui.ColorBlack

	h.view.Clear()
	for i := range menu {
		fmt.Fprintln(h.view, fmt.Sprintf(" %-9s", menu[i]))
	}
	_, err = g.SetCurrentView(MENU)
	return err
}

func NewMenu() *Menu { return &Menu{} }
