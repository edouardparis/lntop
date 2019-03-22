package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

const (
	HEADER = "myheader"
)

type Header struct {
	alias   string
	kind    string
	version string
}

func (h *Header) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(HEADER, x0, y0, x1, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Frame = false
	fmt.Fprintln(v, fmt.Sprintf("[%s %s %s]", h.alias, h.kind, h.version))
	return nil
}

func (h *Header) Update(alias, kind, version string) {
	h.alias = alias
	h.kind = kind
	h.version = version
}

func NewHeader() *Header {
	return &Header{}
}
