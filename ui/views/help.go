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
}

func (h Help) Name() string {
	return HELP
}

func (h Help) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(HELP, x0-1, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Frame = false
	fmt.Fprintln(v, "lntop 0.0.1 - (C) 2019 Edouard Paris")
	fmt.Fprintln(v, "Released under the MIT License")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, fmt.Sprintf("%5s %s",
		color.Cyan("F1 h:"), "show this help screen"))
	_, err = g.SetCurrentView(HELP)
	return err
}

func NewHelp() *Help { return &Help{} }
