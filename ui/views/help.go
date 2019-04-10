package views

import (
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/ui/color"
)

const (
	HELP = "help"
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

func (h *Help) CursorDown() error {
	if h.view != nil {
		cx, cy := h.view.Cursor()
		err := h.view.SetCursor(cx, cy+1)
		if err != nil {
			ox, oy := h.view.Origin()
			err := h.view.SetOrigin(ox, oy+1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Help) CursorUp() error {
	if h.view != nil {
		ox, oy := h.view.Origin()
		cx, cy := h.view.Cursor()
		err := h.view.SetCursor(cx, cy-1)
		if err != nil && oy > 0 {
			err := h.view.SetOrigin(ox, oy-1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Help) CursorRight() error {
	if h.view != nil {
		cx, cy := h.view.Cursor()
		err := h.view.SetCursor(cx+2, cy)
		if err != nil {
			ox, oy := h.view.Origin()
			err := h.view.SetOrigin(ox+2, oy)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *Help) CursorLeft() error {
	if h.view != nil {
		ox, oy := h.view.Origin()
		cx, cy := h.view.Cursor()
		err := h.view.SetCursor(cx-2, cy)
		if err != nil && ox > 1 {
			err := h.view.SetOrigin(ox-2, oy)
			if err != nil {
				return err
			}
		}
	}
	return nil
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
	fmt.Fprintln(h.view, "lntop 0.0.1 - (C) 2019 Edouard Paris")
	fmt.Fprintln(h.view, "Released under the MIT License")
	fmt.Fprintln(h.view, "")
	fmt.Fprintln(h.view, fmt.Sprintf("%5s %s",
		color.Cyan("F1 h:"), "show this help screen"))
	_, err = g.SetCurrentView(HELP)
	return err
}

func NewHelp() *Help { return &Help{} }
