package views

import (
	"bytes"
	"fmt"

	"github.com/jroimartin/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	netmodels "github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
)

const (
	TRANSACTIONS         = "transactions"
	TRANSACTIONS_COLUMNS = "transactions_columns"
	TRANSACTIONS_FOOTER  = "transactions_footer"
)

var DefaultTransactionsColumns = []string{
	"TIME",
	"HEIGHT",
	"CONFIR",
	"AMOUNT",
	"FEE",
	"ADDRESSES",
}

type Transactions struct {
	columns      []transactionsColumn
	columnsView  *gocui.View
	view         *gocui.View
	transactions *models.Transactions

	ox, oy int
	cx, cy int
}

type transactionsColumn struct {
	name    string
	width   int
	display func(*netmodels.Transaction, ...color.Option) string
}

func (c Transactions) Index() int {
	_, oy := c.view.Origin()
	_, cy := c.view.Cursor()
	return cy + oy
}

func (c Transactions) Name() string {
	return TRANSACTIONS
}

func (c *Transactions) Wrap(v *gocui.View) View {
	c.view = v
	return c
}

func (c Transactions) currentColumnIndex() int {
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

func (c Transactions) Origin() (int, int) {
	return c.ox, c.oy
}

func (c Transactions) Cursor() (int, int) {
	return c.cx, c.cy
}

func (c *Transactions) SetCursor(cx, cy int) error {
	err := c.columnsView.SetCursor(cx, 0)
	if err != nil {
		return err
	}

	err = c.view.SetCursor(cx, cy)
	if err != nil {
		return err
	}

	c.cx, c.cy = cx, cy
	return nil
}

func (c *Transactions) SetOrigin(ox, oy int) error {
	err := c.columnsView.SetOrigin(ox, 0)
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

func (c *Transactions) Speed() (int, int, int, int) {
	current := c.currentColumnIndex()
	if current > len(c.columns)-1 {
		return 0, c.columns[current-1].width + 1, 1, 1
	}
	if current == 0 {
		return c.columns[0].width + 1, 0, 1, 1
	}
	return c.columns[current].width + 1,
		c.columns[current-1].width + 1,
		1, 1
}

func (c Transactions) Delete(g *gocui.Gui) error {
	err := g.DeleteView(TRANSACTIONS_COLUMNS)
	if err != nil {
		return err
	}

	err = g.DeleteView(TRANSACTIONS)
	if err != nil {
		return err
	}

	return g.DeleteView(TRANSACTIONS_FOOTER)
}

func (c *Transactions) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	setCursor := false
	c.columnsView, err = g.SetView(TRANSACTIONS_COLUMNS, x0-1, y0, x1+2, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setCursor = true
	}
	c.columnsView.Frame = false
	c.columnsView.BgColor = gocui.ColorGreen
	c.columnsView.FgColor = gocui.ColorBlack

	c.view, err = g.SetView(TRANSACTIONS, x0-1, y0+1, x1+2, y1-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setCursor = true
	}
	c.view.Frame = false
	c.view.Autoscroll = false
	c.view.SelBgColor = gocui.ColorCyan
	c.view.SelFgColor = gocui.ColorBlack
	c.view.Highlight = true
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

	c.display()

	footer, err := g.SetView(TRANSACTIONS_FOOTER, x0-1, y1-2, x1+2, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	footer.Frame = false
	footer.BgColor = gocui.ColorCyan
	footer.FgColor = gocui.ColorBlack
	footer.Clear()
	blackBg := color.Black(color.Background)
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s %s%s %s%s",
		blackBg("F1"), "Help",
		blackBg("F2"), "Menu",
		blackBg("Enter"), "Transaction",
		blackBg("F10"), "Quit",
	))
	return nil
}

func (c *Transactions) display() {
	c.columnsView.Clear()
	var buffer bytes.Buffer
	current := c.currentColumnIndex()
	for i := range c.columns {
		if current == i {
			buffer.WriteString(color.Cyan(color.Background)(c.columns[i].name))
			buffer.WriteString(" ")
			continue
		}
		buffer.WriteString(c.columns[i].name)
		buffer.WriteString(" ")
	}
	fmt.Fprintln(c.columnsView, buffer.String())

	c.view.Clear()
	for _, item := range c.transactions.List() {
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

func NewTransactions(txs *models.Transactions) *Transactions {
	transactions := &Transactions{
		transactions: txs,
	}

	printer := message.NewPrinter(language.English)

	columns := DefaultTransactionsColumns
	transactions.columns = make([]transactionsColumn, len(columns))

	for i := range columns {
		switch columns[i] {
		case "TIME":
			transactions.columns[i] = transactionsColumn{
				name:  fmt.Sprintf("%-15s", columns[i]),
				width: 15,
				display: func(tx *netmodels.Transaction, opts ...color.Option) string {
					return color.Cyan(opts...)(
						fmt.Sprintf("%15s", tx.Date.Format("15:04:05 Jan _2")),
					)
				},
			}
		case "HEIGHT":
			transactions.columns[i] = transactionsColumn{
				name:  fmt.Sprintf("%8s", columns[i]),
				width: 8,
				display: func(tx *netmodels.Transaction, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%8d", tx.BlockHeight))
				},
			}
		case "ADDRESSES":
			transactions.columns[i] = transactionsColumn{
				name:  fmt.Sprintf("%10s", columns[i]),
				width: 10,
				display: func(tx *netmodels.Transaction, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%10d", len(tx.DestAddresses)))
				},
			}
		case "FEE":
			transactions.columns[i] = transactionsColumn{
				name:  fmt.Sprintf("%8s", columns[i]),
				width: 8,
				display: func(tx *netmodels.Transaction, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%8d", tx.TotalFees))
				},
			}
		case "CONFIR":
			transactions.columns[i] = transactionsColumn{
				name:  fmt.Sprintf("%8s", columns[i]),
				width: 8,
				display: func(tx *netmodels.Transaction, opts ...color.Option) string {
					n := fmt.Sprintf("%8d", tx.NumConfirmations)
					if tx.NumConfirmations < 6 {
						return color.Yellow(opts...)(n)
					}
					return color.Green(opts...)(n)
				},
			}
		case "TXHASH":
			transactions.columns[i] = transactionsColumn{
				name:  fmt.Sprintf("%-64s", columns[i]),
				width: 64,
				display: func(tx *netmodels.Transaction, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%13s", tx.TxHash))
				},
			}
		case "BLOCKHASH":
			transactions.columns[i] = transactionsColumn{
				name: fmt.Sprintf("%-64s", columns[i]),
				display: func(tx *netmodels.Transaction, opts ...color.Option) string {
					return color.White(opts...)(fmt.Sprintf("%13s", tx.TxHash))
				},
			}
		case "AMOUNT":
			transactions.columns[i] = transactionsColumn{
				name:  fmt.Sprintf("%13s", columns[i]),
				width: 13,
				display: func(tx *netmodels.Transaction, opts ...color.Option) string {
					return color.White(opts...)(printer.Sprintf("%13d", tx.Amount))
				},
			}
		default:
			transactions.columns[i] = transactionsColumn{
				name:  fmt.Sprintf("%-21s", columns[i]),
				width: 21,
				display: func(tx *netmodels.Transaction, opts ...color.Option) string {
					return "column does not exist"
				},
			}
		}

	}
	return transactions
}
