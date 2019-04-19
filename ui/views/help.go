package views

import (
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/ui/color"
)

const (
	version = "v0.0.2"
	HELP    = "help"
)

type Help struct {
	view *gocui.View
}

func (h Help) Name() string {
	return HELP
}

func (h *Help) Wrap(v *gocui.View) view {
	h.view = v
	return h
}

func (h Help) Delete(g *gocui.Gui) error {
	return g.DeleteView(HELP)
}

func (h *Help) CursorDown() error {
	return cursorDown(h.view, 1)
}

func (h *Help) CursorUp() error {
	return cursorUp(h.view, 1)
}

func (h *Help) CursorRight() error {
	return cursorRight(h.view, 1)
}

func (h *Help) CursorLeft() error {
	return cursorLeft(h.view, 1)
}

func (h Help) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	h.view, err = g.SetView(HELP, x0-1, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	h.view.Frame = false
	fmt.Fprintln(h.view, fmt.Sprintf("lntop %s - (C) 2019 Edouard Paris", version))
	fmt.Fprintln(h.view, "Released under the MIT License")
	fmt.Fprintln(h.view, "")
	fmt.Fprintln(h.view, fmt.Sprintf("%6s %s",
		color.Cyan("F1  h:"), "show/close this help screen"))
	fmt.Fprintln(h.view, fmt.Sprintf("%6s %s",
		color.Cyan("F10 q:"), "quit"))
	_, err = g.SetCurrentView(HELP)
	return err
}

func NewHelp() *Help { return &Help{} }
