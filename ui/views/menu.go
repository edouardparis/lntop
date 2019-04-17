package views

import (
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/ui/color"
)

const (
	MENU = "menu"
)

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

func (h Menu) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	h.view, err = g.SetView(MENU, x0-1, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	h.view.Frame = false
	fmt.Fprintln(h.view, fmt.Sprintf("lntop %s - (C) 2019 Edouard Paris", version))
	fmt.Fprintln(h.view, "Released under the MIT License")
	fmt.Fprintln(h.view, "")
	fmt.Fprintln(h.view, fmt.Sprintf("%5s %s",
		color.Cyan("F1 h:"), "show this menu screen"))
	_, err = g.SetCurrentView(MENU)
	return err
}

func NewMenu() *Menu { return &Menu{} }
