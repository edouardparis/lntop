package ui

import (
	"github.com/jroimartin/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func setKeyBinding(c *controller, g *gocui.Gui) error {
	err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyF10, gocui.ModNone, quit)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", 'q', gocui.ModNone, quit)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, c.cursorUp)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, c.cursorDown)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, c.cursorLeft)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, c.cursorRight)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, c.OnEnter)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyF1, gocui.ModNone, c.Help)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", 'h', gocui.ModNone, c.Help)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", gocui.KeyF2, gocui.ModNone, c.Menu)
	if err != nil {
		return err
	}

	err = g.SetKeybinding("", 'm', gocui.ModNone, c.Menu)
	if err != nil {
		return err
	}

	return nil
}
