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
	col          int
	columns      []transactionsColumn
	columnsView  *gocui.View
	view         *gocui.View
	transactions *models.Transactions
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

func (c *Transactions) Wrap(v *gocui.View) view {
	c.view = v
	return c
}

func (c *Transactions) CursorDown() error {
	return cursorDown(c.view, 1)
}

func (c *Transactions) CursorUp() error {
	return cursorUp(c.view, 1)
}

func (c *Transactions) CursorRight() error {
	if c.col > len(c.columns)-1 {
		return nil
	}
	speed := c.columns[c.col].width + 1
	if c.col == len(c.columns)-1 {
		speed := c.columns[c.col].width + 1
		err := cursorRight(c.columnsView, speed)
		if err != nil {
			return err
		}

		err = cursorRight(c.view, speed)
		if err != nil {
			return err
		}
		err = cursorLeft(c.columnsView, speed)
		if err != nil {
			return err
		}

		return cursorLeft(c.view, speed)
	}
	c.col++
	err := cursorRight(c.columnsView, speed)
	if err != nil {
		return err
	}

	return cursorRight(c.view, speed)
}

func (c *Transactions) CursorLeft() error {
	if c.col == 0 {
		return nil
	}
	speed := c.columns[c.col-1].width + 1
	c.col--
	err := cursorLeft(c.columnsView, speed)
	if err != nil {
		return err
	}

	return cursorLeft(c.view, speed)
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
	c.columnsView, err = g.SetView(TRANSACTIONS_COLUMNS, x0-1, y0, x1+2, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	c.columnsView.Frame = false
	c.columnsView.BgColor = gocui.ColorGreen
	c.columnsView.FgColor = gocui.ColorBlack

	c.view, err = g.SetView(TRANSACTIONS, x0-1, y0+1, x1+2, y1-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	c.view.Frame = false
	c.view.Autoscroll = false
	c.view.SelBgColor = gocui.ColorCyan
	c.view.SelFgColor = gocui.ColorBlack
	c.view.Highlight = true

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
	for i := range c.columns {
		if c.col == i {
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
			if c.col == i {
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
