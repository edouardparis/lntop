package views

import (
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop/ui/color"
)

const (
	version = "v0.1.0"
	HELP    = "help"
)

type Help struct {
	view *gocui.View
}

func (h Help) Name() string {
	return HELP
}

func (h *Help) Wrap(v *gocui.View) View {
	h.view = v
	return h
}

func (h Help) Delete(g *gocui.Gui) error {
	return g.DeleteView(HELP)
}

func (h Help) Origin() (int, int) {
	return h.view.Origin()
}

func (h Help) Cursor() (int, int) {
	return h.view.Cursor()
}

func (h Help) Speed() (int, int, int, int) {
	return 1, 1, 1, 1
}

func (h Help) Limits() (pageSize int, fullSize int) {
	return 0, 0
}

func (h *Help) SetCursor(x, y int) error {
	return h.view.SetCursor(x, y)
}

func (h *Help) SetOrigin(x, y int) error {
	return h.view.SetOrigin(x, y)
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
	cyan := color.Cyan()
	fmt.Fprintln(h.view, fmt.Sprintf("lntop %s - (C) 2019 Edouard Paris", version))
	fmt.Fprintln(h.view, "Released under the MIT License")
	fmt.Fprintln(h.view, "")
	fmt.Fprintln(h.view, fmt.Sprintf("%6s %s",
		cyan("F1  h:"), "show/close this help screen"))
	fmt.Fprintln(h.view, fmt.Sprintf("%6s %s",
		cyan("F2  m:"), "show/close the menu sidebar"))
	fmt.Fprintln(h.view, fmt.Sprintf("%6s %s",
		cyan("F10 q:"), "quit"))

	fmt.Fprintln(h.view, "")
	fmt.Fprintln(h.view, fmt.Sprintf("%6s %s",
		cyan("  a d:"), "apply asc/desc order to the rows according to the selected column value"))
	_, err = g.SetCurrentView(HELP)
	return err
}

func NewHelp() *Help { return &Help{} }
