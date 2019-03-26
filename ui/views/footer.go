package views

import "github.com/jroimartin/gocui"

const (
	FOOTER = "footer"
)

type Footer struct{}

func (f *Footer) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(FOOTER, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Frame = false
	v.BgColor = gocui.ColorCyan
	v.FgColor = gocui.ColorBlack | gocui.AttrBold
	return nil
}

func NewFooter() *Footer {
	return &Footer{}
}
