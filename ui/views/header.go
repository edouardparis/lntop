package views

import (
	"fmt"
	"regexp"

	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
	"github.com/jroimartin/gocui"
)

const (
	HEADER = "myheader"
)

var versionReg = regexp.MustCompile(`(\d+\.)?(\d+\.)?(\*|\d+)`)

type Header struct {
	Info *models.Info
}

func (h *Header) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(HEADER, x0, y0, x1, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Frame = false

	version := h.Info.Version
	matches := versionReg.FindStringSubmatch(h.Info.Version)
	if len(matches) > 0 {
		version = matches[0]
	}

	fmt.Fprintln(v, fmt.Sprintf("%s %s %s %s",
		color.CyanBg(h.Info.Alias),
		color.Cyan(fmt.Sprintf("%s-v%s", "lnd", version)),
		fmt.Sprintf("%s %d", color.Cyan("height:"), h.Info.BlockHeight),
		fmt.Sprintf("%s %d", color.Cyan("peers:"), h.Info.NumPeers),
	))
	return nil
}

func NewHeader(info *models.Info) *Header {
	return &Header{Info: info}
}
