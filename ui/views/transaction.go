package views

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
)

const (
	TRANSACTION        = "transaction"
	TRANSACTION_HEADER = "transaction_header"
	TRANSACTION_FOOTER = "transaction_footer"
)

type Transaction struct {
	view        *gocui.View
	transaction *models.Transaction
}

func (c Transaction) Name() string {
	return TRANSACTION
}

func (c Transaction) Empty() bool {
	return c.transaction == nil
}

func (c *Transaction) Wrap(v *gocui.View) view {
	c.view = v
	return c
}

func (c *Transaction) CursorDown() error {
	return cursorDown(c.view, 1)
}

func (c *Transaction) CursorUp() error {
	return cursorUp(c.view, 1)
}

func (c *Transaction) CursorRight() error {
	return cursorRight(c.view, 1)
}

func (c *Transaction) CursorLeft() error {
	return cursorLeft(c.view, 1)
}

func (c *Transaction) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	header, err := g.SetView(TRANSACTION_HEADER, x0-1, y0, x1+2, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	header.Frame = false
	header.BgColor = gocui.ColorGreen
	header.FgColor = gocui.ColorBlack | gocui.AttrBold
	header.Clear()
	fmt.Fprintln(header, "Transaction")

	v, err := g.SetView(TRANSACTION, x0-1, y0+1, x1+2, y1-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	v.Frame = false
	c.view = v
	c.display()

	footer, err := g.SetView(TRANSACTION_FOOTER, x0-1, y1-2, x1, y1)
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
		color.BlackBg("Enter"), "Transactions",
		color.BlackBg("F10"), "Quit",
	))
	return nil
}

func (c Transaction) Delete(g *gocui.Gui) error {
	err := g.DeleteView(TRANSACTION_HEADER)
	if err != nil {
		return err
	}

	err = g.DeleteView(TRANSACTION)
	if err != nil {
		return err
	}

	return g.DeleteView(TRANSACTION_FOOTER)
}

func (c *Transaction) display() {
	p := message.NewPrinter(language.English)
	v := c.view
	v.Clear()
	transaction := c.transaction.Item
	fmt.Fprintln(v, color.Green(" [ Transaction ]"))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		color.Cyan("           Date:"), transaction.Date.Format("15:04:05 Jan _2")))
	fmt.Fprintln(v, p.Sprintf("%s %d",
		color.Cyan("         Amount:"), transaction.Amount))
	fmt.Fprintln(v, p.Sprintf("%s %d",
		color.Cyan("            Fee:"), transaction.TotalFees))
	fmt.Fprintln(v, p.Sprintf("%s %d",
		color.Cyan("    BlockHeight:"), transaction.BlockHeight))
	fmt.Fprintln(v, p.Sprintf("%s %d",
		color.Cyan("NumConfirmations:"), transaction.NumConfirmations))
	fmt.Fprintln(v, p.Sprintf("%s %s",
		color.Cyan("       BlockHash:"), transaction.BlockHash))
	fmt.Fprintln(v, fmt.Sprintf("%s %s",
		color.Cyan("         TxHash:"), transaction.TxHash))
	fmt.Fprintln(v, "[addresses]")
}

func NewTransaction(transaction *models.Transaction) *Transaction {
	return &Transaction{transaction: transaction}
}
