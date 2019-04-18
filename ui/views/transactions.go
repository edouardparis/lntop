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
	"TXHASH",
	"AMOUNT",
}

type Transactions struct {
	columns      []transactionsColumn
	columnsView  *gocui.View
	view         *gocui.View
	transactions *models.Transactions
}

type transactionsColumn struct {
	name    string
	display func(*netmodels.Transaction) string
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
	err := cursorRight(c.columnsView, 2)
	if err != nil {
		return err
	}

	return cursorRight(c.view, 2)
}

func (c *Transactions) CursorLeft() error {
	err := cursorLeft(c.columnsView, 2)
	if err != nil {
		return err
	}

	return cursorLeft(c.view, 2)
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
		_, err = g.SetCurrentView(TRANSACTIONS)
		if err != nil {
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
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s %s%s %s%s",
		color.BlackBg("F1"), "Help",
		color.BlackBg("F2"), "Menu",
		color.BlackBg("Enter"), "Transaction",
		color.BlackBg("F10"), "Quit",
	))
	return nil
}

func (c *Transactions) display() {
	c.columnsView.Clear()
	var buffer bytes.Buffer
	for i := range c.columns {
		buffer.WriteString(c.columns[i].name)
		buffer.WriteString(" ")
	}
	fmt.Fprintln(c.columnsView, buffer.String())

	c.view.Clear()
	for _, item := range c.transactions.List() {
		var buffer bytes.Buffer
		for i := range c.columns {
			buffer.WriteString(c.columns[i].display(item))
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
		case "TXHASH":
			transactions.columns[i] = transactionsColumn{
				name: fmt.Sprintf("%-13s", columns[i]),
				display: func(tx *netmodels.Transaction) string {
					return fmt.Sprintf("%13s", tx.TxHash)
				},
			}
		case "AMOUNT":
			transactions.columns[i] = transactionsColumn{
				name: fmt.Sprintf("%-13s", columns[i]),
				display: func(tx *netmodels.Transaction) string {
					return printer.Sprintf("%13d", tx.Amount)
				},
			}
		default:
			transactions.columns[i] = transactionsColumn{
				name: fmt.Sprintf("%-21s", columns[i]),
				display: func(tx *netmodels.Transaction) string {
					return "column does not exist"
				},
			}
		}

	}
	return transactions
}
