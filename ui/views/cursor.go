package views

import "github.com/jroimartin/gocui"

func cursorDown(v *gocui.View, speed int) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	err := v.SetCursor(cx, cy+speed)
	if err != nil {
		ox, oy := v.Origin()
		err := v.SetOrigin(ox, oy+speed)
		if err != nil {
			return err
		}
	}
	return nil
}

func cursorUp(v *gocui.View, speed int) error {
	if v == nil {
		return nil
	}
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	err := v.SetCursor(cx, cy-speed)
	if err != nil && oy >= speed {
		err := v.SetOrigin(ox, oy-speed)
		if err != nil {
			return err
		}
	}
	return nil
}

func cursorRight(v *gocui.View, speed int) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	err := v.SetCursor(cx+speed, cy)
	if err != nil {
		ox, oy := v.Origin()
		err := v.SetOrigin(ox+speed, oy)
		if err != nil {
			return err
		}
	}
	return nil
}

func cursorLeft(v *gocui.View, speed int) error {
	if v == nil {
		return nil
	}
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	err := v.SetCursor(cx-speed, cy)
	if err != nil && ox >= speed {
		err := v.SetOrigin(ox-speed, oy)
		if err != nil {
			return err
		}
	}
	return nil
}
