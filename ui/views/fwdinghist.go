package views

import (
	"bytes"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/edouardparis/lntop/config"
	netmodels "github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
)

const (
	FWDINGHIST         = "fwdinghist"
	FWDINGHIST_COLUMNS = "fwdinghist_columns"
	FWDINGHIST_FOOTER  = "fwdinghist_footer"
)

var DefaultFwdinghistColumns = []string{
	"ALIAS_IN",
	"ALIAS_OUT",
	"AMT_IN",
	"AMT_OUT",
	"FEE",
	"TIMESTAMP_NS",
	"CHAN_ID_IN",
	"CHAN_ID_OUT",
}

type FwdingHist struct {
	cfg *config.View

	columns           []fwdinghistColumn
	columnHeadersView *gocui.View
	view              *gocui.View
	fwdinghist        *models.FwdingHist

	ox, oy int
	cx, cy int
}

type fwdinghistColumn struct {
	name    string
	width   int
	sorted  bool
	sort    func(models.Order) models.FwdinghistSort
	display func(*netmodels.ForwardingEvent, ...color.Option) string
}

func (c FwdingHist) Index() int {
	_, oy := c.view.Origin()
	_, cy := c.view.Cursor()
	return cy + oy
}

func (c FwdingHist) Name() string {
	return FWDINGHIST
}

func (c *FwdingHist) Wrap(v *gocui.View) View {
	c.view = v
	return c
}

func (c FwdingHist) currentColumnIndex() int {
	x := c.ox + c.cx
	index := 0
	sum := 0
	for i := range c.columns {
		sum += c.columns[i].width + 1
		if x < sum {
			return index
		}
		index++
	}
	return index
}

func (c FwdingHist) Origin() (int, int) {
	return c.ox, c.oy
}

func (c FwdingHist) Cursor() (int, int) {
	return c.cx, c.cy
}

func (c *FwdingHist) SetCursor(cx, cy int) error {
	if err := cursorCompat(c.columnHeadersView, cx, 0); err != nil {
		return err
	}
	err := c.columnHeadersView.SetCursor(cx, 0)
	if err != nil {
		return err
	}

	if err := cursorCompat(c.view, cx, cy); err != nil {
		return err
	}
	err = c.view.SetCursor(cx, cy)
	if err != nil {
		return err
	}

	c.cx, c.cy = cx, cy
	return nil
}

func (c *FwdingHist) SetOrigin(ox, oy int) error {
	err := c.columnHeadersView.SetOrigin(ox, 0)
	if err != nil {
		return err
	}
	err = c.view.SetOrigin(ox, oy)
	if err != nil {
		return err
	}

	c.ox, c.oy = ox, oy
	return nil
}

func (c *FwdingHist) Speed() (int, int, int, int) {
	current := c.currentColumnIndex()
	up := 0
	down := 0
	if c.Index() > 0 {
		up = 1
	}
	if c.Index() < c.fwdinghist.Len()-1 {
		down = 1
	}
	if current > len(c.columns)-1 {
		return 0, c.columns[current-1].width + 1, down, up
	}
	if current == 0 {
		return c.columns[0].width + 1, 0, down, up
	}
	return c.columns[current].width + 1,
		c.columns[current-1].width + 1,
		down, up
}

func (c *FwdingHist) Limits() (pageSize int, fullSize int) {
	_, pageSize = c.view.Size()
	fullSize = c.fwdinghist.Len()
	return
}

func (c *FwdingHist) Sort(column string, order models.Order) {
	if column == "" {
		index := c.currentColumnIndex()
		if index >= len(c.columns) {
			return
		}
		col := c.columns[index]
		if col.sort == nil {
			return
		}

		c.fwdinghist.Sort(col.sort(order))
		for i := range c.columns {
			c.columns[i].sorted = (i == index)
		}
	}
}

func (c FwdingHist) Delete(g *gocui.Gui) error {
	err := g.DeleteView(FWDINGHIST_COLUMNS)
	if err != nil {
		return err
	}

	err = g.DeleteView(FWDINGHIST)
	if err != nil {
		return err
	}

	return g.DeleteView(FWDINGHIST_FOOTER)
}

func (c *FwdingHist) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	setCursor := false
	c.columnHeadersView, err = g.SetView(FWDINGHIST_COLUMNS, x0-1, y0, x1+2, y0+2, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setCursor = true
	}
	c.columnHeadersView.Frame = false
	c.columnHeadersView.BgColor = gocui.ColorGreen
	c.columnHeadersView.FgColor = gocui.ColorBlack

	c.view, err = g.SetView(FWDINGHIST, x0-1, y0+1, x1+2, y1-1, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setCursor = true
	}
	c.view.Frame = false
	c.view.Autoscroll = false
	c.view.SelBgColor = gocui.ColorCyan
	c.view.SelFgColor = gocui.ColorBlack | gocui.AttrDim
	c.view.Highlight = true
	c.display()

	if setCursor {
		ox, oy := c.Origin()
		err := c.SetOrigin(ox, oy)
		if err != nil {
			return err
		}

		cx, cy := c.Cursor()
		err = c.SetCursor(cx, cy)
		if err != nil {
			return err
		}
	}

	footer, err := g.SetView(FWDINGHIST_FOOTER, x0-1, y1-2, x1+2, y1, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	footer.Frame = false
	footer.BgColor = gocui.ColorCyan
	footer.FgColor = gocui.ColorBlack
	footer.Rewind()
	blackBg := color.Black(color.Background)
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s %s%s",
		blackBg("F2"), "Menu",
		blackBg("Enter"), "FwdingHist",
		blackBg("F10"), "Quit",
	))
	return nil
}

func (c *FwdingHist) display() {
	c.columnHeadersView.Rewind()
	var buffer bytes.Buffer
	current := c.currentColumnIndex()
	for i := range c.columns {
		if current == i {
			buffer.WriteString(color.Cyan(color.Background)(c.columns[i].name))
			buffer.WriteString(" ")
			continue
		} else if c.columns[i].sorted {
			buffer.WriteString(color.Magenta(color.Background)(c.columns[i].name))
			buffer.WriteString(" ")
			continue
		}
		buffer.WriteString(c.columns[i].name)
		buffer.WriteString(" ")
	}
	fmt.Fprintln(c.columnHeadersView, buffer.String())

	c.view.Rewind()
	for _, item := range c.fwdinghist.List() {
		var buffer bytes.Buffer
		for i := range c.columns {
			var opt color.Option
			if current == i {
				opt = color.Bold
			}
			buffer.WriteString(c.columns[i].display(item, opt))
			buffer.WriteString(" ")
		}
		fmt.Fprintln(c.view, buffer.String())
	}
}

func NewFwdingHist(cfg *config.View, hist *models.FwdingHist) *FwdingHist {
	fwdinghist := &FwdingHist{
		cfg:        cfg,
		fwdinghist: hist,
	}

	printer := message.NewPrinter(language.English)

	columns := DefaultFwdinghistColumns
	if cfg != nil && len(cfg.Columns) != 0 {
		columns = cfg.Columns
	}

	fwdinghist.columns = make([]fwdinghistColumn, len(columns))

	for i := range columns {
		switch columns[i] {
		case "ALIAS_IN":
			fwdinghist.columns[i] = fwdinghistColumn{
				width: 30,
				name:  fmt.Sprintf("%30s", columns[i]),
				sort: func(order models.Order) models.FwdinghistSort {
					return func(e1, e2 *netmodels.ForwardingEvent) bool {
						return models.StringSort(e1.PeerAliasIn, e2.PeerAliasOut, order)
					}
				},
				display: func(e *netmodels.ForwardingEvent, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%30s", e.PeerAliasIn))
				},
			}
		case "ALIAS_OUT":
			fwdinghist.columns[i] = fwdinghistColumn{
				width: 30,
				name:  fmt.Sprintf("%30s", columns[i]),
				sort: func(order models.Order) models.FwdinghistSort {
					return func(e1, e2 *netmodels.ForwardingEvent) bool {
						return models.StringSort(e1.PeerAliasOut, e2.PeerAliasOut, order)
					}
				},
				display: func(e *netmodels.ForwardingEvent, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%30s", e.PeerAliasOut))
				},
			}
		case "CHAN_ID_IN":
			fwdinghist.columns[i] = fwdinghistColumn{
				width: 19,
				name:  fmt.Sprintf("%19s", columns[i]),
				sort: func(order models.Order) models.FwdinghistSort {
					return func(e1, e2 *netmodels.ForwardingEvent) bool {
						return models.UInt64Sort(e1.ChanIdIn, e2.ChanIdIn, order)
					}
				},
				display: func(e *netmodels.ForwardingEvent, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%19d", e.ChanIdIn))
				},
			}
		case "CHAN_ID_OUT":
			fwdinghist.columns[i] = fwdinghistColumn{
				width: 19,
				name:  fmt.Sprintf("%19s", columns[i]),
				sort: func(order models.Order) models.FwdinghistSort {
					return func(e1, e2 *netmodels.ForwardingEvent) bool {
						return models.UInt64Sort(e1.ChanIdOut, e2.ChanIdOut, order)
					}
				},
				display: func(e *netmodels.ForwardingEvent, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%19d", e.ChanIdOut))
				},
			}
		case "AMT_IN":
			fwdinghist.columns[i] = fwdinghistColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", "RECEIVED"),
				sort: func(order models.Order) models.FwdinghistSort {
					return func(e1, e2 *netmodels.ForwardingEvent) bool {
						return models.UInt64Sort(e1.AmtIn, e2.AmtIn, order)
					}
				},
				display: func(e *netmodels.ForwardingEvent, opts ...color.Option) string {
					return color.White(opts...)(printer.Sprintf("%12d", e.AmtIn))
				},
			}
		case "AMT_OUT":
			fwdinghist.columns[i] = fwdinghistColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", "SENT"),
				sort: func(order models.Order) models.FwdinghistSort {
					return func(e1, e2 *netmodels.ForwardingEvent) bool {
						return models.UInt64Sort(e1.AmtOut, e2.AmtOut, order)
					}
				},
				display: func(e *netmodels.ForwardingEvent, opts ...color.Option) string {
					return color.White(opts...)(printer.Sprintf("%12d", e.AmtOut))
				},
			}
		case "FEE":
			fwdinghist.columns[i] = fwdinghistColumn{
				name:  fmt.Sprintf("%9s", "EARNED"),
				width: 9,
				sort: func(order models.Order) models.FwdinghistSort {
					return func(e1, e2 *netmodels.ForwardingEvent) bool {
						return models.UInt64Sort(e1.Fee, e2.Fee, order)
					}
				},
				display: func(e *netmodels.ForwardingEvent, opts ...color.Option) string {
					return fee(e.Fee)
				},
			}
		case "TIMESTAMP_NS":
			fwdinghist.columns[i] = fwdinghistColumn{
				name:  fmt.Sprintf("%15s", "TIME"),
				width: 20,
				display: func(e *netmodels.ForwardingEvent, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%20s", e.EventTime.Format("15:04:05 Jan _2")))
				},
			}
		default:
			fwdinghist.columns[i] = fwdinghistColumn{
				name:  fmt.Sprintf("%-21s", columns[i]),
				width: 21,
				display: func(tx *netmodels.ForwardingEvent, opts ...color.Option) string {
					return "column does not exist"
				},
			}
		}

	}
	return fwdinghist
}

func fee(fee uint64, opts ...color.Option) string {
	if fee >= 0 && fee < 100 {
		return color.Cyan(opts...)(fmt.Sprintf("%9d", fee))
	} else if fee >= 100 && fee < 999 {
		return color.Green(opts...)(fmt.Sprintf("%9d", fee))
	}

	return color.Yellow(opts...)(fmt.Sprintf("%9d", fee))
}
