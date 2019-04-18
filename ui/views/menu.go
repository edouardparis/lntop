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
	"CHAN",
	"TX",
}

type Menu struct {
	view *gocui.View
}

func (h Menu) Name() string {
	return MENU
}

func (h *Menu) Wrap(v *gocui.View) view {
	h.view = v
	return h
}

func (h *Menu) CursorDown() error {
	return cursorDown(h.view, 1)
}

func (h *Menu) CursorUp() error {
	return cursorUp(h.view, 1)
}

func (h *Menu) CursorRight() error {
	return nil
}

func (h *Menu) CursorLeft() error {
	return nil
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

	h.view, err = g.SetView(MENU, x0, y0+1, x1, y1-2)
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
		fmt.Fprintln(h.view, menu[i])
	}
	_, err = g.SetCurrentView(MENU)
	return err
}

func NewMenu() *Menu { return &Menu{} }
