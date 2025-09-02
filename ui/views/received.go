package views

import (
	"bytes"
	"fmt"
	"time"

	"github.com/awesome-gocui/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/edouardparis/lntop/config"
	netmodels "github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
)

const (
	RECEIVED         = "received"
	RECEIVED_COLUMNS = "received_columns"
	RECEIVED_FOOTER  = "received_footer"
)

var DefaultReceivedColumns = []string{
	"TYPE",
	"TIME",
	"AMOUNT",
	"MEMO",
	"R_HASH",
}

type Received struct {
	cfg *config.View

	columns           []receivedColumn
	columnHeadersView *gocui.View
	view              *gocui.View
	received          *models.Received

	ox, oy int
	cx, cy int
}

type receivedColumn struct {
	name    string
	width   int
	display func(*netmodels.Invoice, ...color.Option) string
}

func (c Received) Name() string             { return RECEIVED }
func (c *Received) Wrap(v *gocui.View) View { c.view = v; return c }
func (c Received) Origin() (int, int)       { return c.ox, c.oy }
func (c Received) Cursor() (int, int)       { return c.cx, c.cy }

func (c *Received) SetCursor(cx, cy int) error {
	if err := cursorCompat(c.columnHeadersView, cx, 0); err != nil {
		return err
	}
	if err := c.columnHeadersView.SetCursor(cx, 0); err != nil {
		return err
	}
	if err := cursorCompat(c.view, cx, cy); err != nil {
		return err
	}
	if err := c.view.SetCursor(cx, cy); err != nil {
		return err
	}
	c.cx, c.cy = cx, cy
	return nil
}

func (c *Received) SetOrigin(ox, oy int) error {
	if err := c.columnHeadersView.SetOrigin(ox, 0); err != nil {
		return err
	}
	if err := c.view.SetOrigin(ox, oy); err != nil {
		return err
	}
	c.ox, c.oy = ox, oy
	return nil
}

func (c *Received) Speed() (int, int, int, int) {
	up, down := 0, 0
	if c.Index() > 0 {
		up = 1
	}
	if c.Index() < c.received.Len()-1 {
		down = 1
	}
	// columns have static widths: TYPE(7) TIME(25) AMOUNT(12) MEMO(40) R_HASH(64)
	return 7, 0, down, up
}

func (c *Received) Limits() (int, int) {
	_, page := c.view.Size()
	full := c.received.Len()
	return page, full
}

func (c Received) Index() int { _, oy := c.view.Origin(); _, cy := c.view.Cursor(); return cy + oy }

func (c Received) Delete(g *gocui.Gui) error {
	if err := g.DeleteView(RECEIVED_COLUMNS); err != nil {
		return err
	}
	if err := g.DeleteView(RECEIVED); err != nil {
		return err
	}
	return g.DeleteView(RECEIVED_FOOTER)
}

func (c *Received) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	setCursor := false
	c.columnHeadersView, err = g.SetView(RECEIVED_COLUMNS, x0-1, y0, x1+2, y0+2, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setCursor = true
	}
	c.columnHeadersView.Frame = false
	c.columnHeadersView.BgColor = gocui.ColorGreen
	c.columnHeadersView.FgColor = gocui.ColorBlack

	c.view, err = g.SetView(RECEIVED, x0-1, y0+1, x1+2, y1-1, 0)
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
		if err := c.SetOrigin(c.ox, c.oy); err != nil {
			return err
		}
		if err := c.SetCursor(c.cx, c.cy); err != nil {
			return err
		}
	}

	footer, err := g.SetView(RECEIVED_FOOTER, x0-1, y1-2, x1+2, y1, 0)
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
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s",
		blackBg("F2"), "Menu",
		blackBg("F10"), "Quit",
	))
	return nil
}

func (c *Received) display() {
	c.columnHeadersView.Rewind()
	var buffer bytes.Buffer
	for i := range c.columns {
		buffer.WriteString(c.columns[i].name)
		buffer.WriteString(" ")
	}
	fmt.Fprintln(c.columnHeadersView, buffer.String())

	c.view.Rewind()
	for _, inv := range c.received.List() {
		var b bytes.Buffer
		for i := range c.columns {
			b.WriteString(c.columns[i].display(inv))
			b.WriteString(" ")
		}
		fmt.Fprintln(c.view, b.String())
	}
}

func NewReceived(cfg *config.View, rec *models.Received) *Received {
	received := &Received{cfg: cfg, received: rec}

	printer := message.NewPrinter(language.English)

	cols := DefaultReceivedColumns
	if cfg != nil && len(cfg.Columns) != 0 {
		cols = cfg.Columns
	}

	received.columns = make([]receivedColumn, len(cols))
	for i := range cols {
		switch cols[i] {
		case "TYPE":
			received.columns[i] = receivedColumn{
				width: 7,
				name:  fmt.Sprintf("%-7s", cols[i]),
				display: func(inv *netmodels.Invoice, opts ...color.Option) string {
					label := "invoice"
					col := color.White(opts...)
					if inv.Kind == netmodels.KindKeysend || inv.PaymentRequest == "" {
						label = "keysend"
						col = color.White(opts...) // make keysend white like invoice
					}
					return col(fmt.Sprintf("%-7s", label))
				},
			}
		case "TIME":
			received.columns[i] = receivedColumn{
				width: 25,
				name:  fmt.Sprintf("%25s", cols[i]),
				display: func(inv *netmodels.Invoice, opts ...color.Option) string {
					// Prefer settle date, fallback to creation
					ts := inv.SettleDate
					if ts == 0 {
						ts = inv.CreationDate
					}
					// Show time with year appended, preserving original style
					return color.White(opts...)(fmt.Sprintf("%25s", time.Unix(ts, 0).Format("15:04:05 Jan _2 2006")))
				},
			}
		case "AMOUNT":
			received.columns[i] = receivedColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", cols[i]),
				display: func(inv *netmodels.Invoice, opts ...color.Option) string {
					amt := inv.AmountPaid
					if amt == 0 {
						amt = inv.Amount
					}
					return color.Yellow(opts...)(printer.Sprintf("%12d", amt))
				},
			}
		case "MEMO":
			received.columns[i] = receivedColumn{
				width: 40,
				name:  fmt.Sprintf("%-40s", cols[i]),
				display: func(inv *netmodels.Invoice, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%-40s", inv.Description))
				},
			}
		case "R_HASH":
			received.columns[i] = receivedColumn{
				width: 64,
				name:  fmt.Sprintf("%-64s", cols[i]),
				display: func(inv *netmodels.Invoice, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%-64s", inv.GetRHash()))
				},
			}
		default:
			received.columns[i] = receivedColumn{
				width:   10,
				name:    fmt.Sprintf("%-10s", cols[i]),
				display: func(inv *netmodels.Invoice, opts ...color.Option) string { return "" },
			}
		}
	}
	return received
}
