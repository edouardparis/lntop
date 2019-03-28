package views

import (
	"fmt"
	"regexp"

	"github.com/edouardparis/lntop/ui/color"
	"github.com/jroimartin/gocui"
)

const (
	HEADER = "myheader"
)

var versionReg = regexp.MustCompile(`(\d+\.)?(\d+\.)?(\*|\d+)`)

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

	version := h.version
	matches := versionReg.FindStringSubmatch(h.version)
	if len(matches) > 0 {
		version = matches[0]
	}

	fmt.Fprintln(v, fmt.Sprintf("%s %s %s",
		color.CyanBg(h.alias), color.Cyan(h.kind), color.Cyan(version)))
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
