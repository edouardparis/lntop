package views

import (
	"fmt"

	"github.com/jroimartin/gocui"

	"github.com/edouardparis/lntop"
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
	_, pageSize = h.view.Size()
	fullSize = len(h.view.BufferLines()) - 1
	return
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
	fmt.Fprintf(h.view, "lntop %s - (C) 2019 Edouard Paris\n", lntop.Version)
	fmt.Fprintln(h.view, "Released under the MIT License")
	fmt.Fprintln(h.view, "")
	fmt.Fprintf(h.view, "%6s show/close this help screen\n", cyan("F1  h:"))
	fmt.Fprintf(h.view, "%6s show/close the menu sidebar\n", cyan("F2  m:"))
	fmt.Fprintf(h.view, "%6s quit\n", cyan("F10 q:"))
	fmt.Fprintln(h.view, "")
	fmt.Fprintf(h.view, "%6s apply asc/desc order to the rows according to the selected column value\n",
		cyan("  a d:"))
	_, err = g.SetCurrentView(HELP)
	return err
}

func NewHelp() *Help { return &Help{} }
