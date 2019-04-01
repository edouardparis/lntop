package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	HELP = "help"
)

func SetHelp(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(HELP, x0-1, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Frame = false
	fmt.Fprintln(v, "HELP")
	_, err = g.SetCurrentView(HELP)
	return err
}
