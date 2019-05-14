package cursor

type View interface {
	Cursor() (int, int)
	Origin() (int, int)
	Speed() (int, int, int, int)
	SetCursor(int, int) error
	SetOrigin(int, int) error
}

func Down(v View) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	_, _, sy, _ := v.Speed()
	err := v.SetCursor(cx, cy+sy)
	if err != nil {
		ox, oy := v.Origin()
		err := v.SetOrigin(ox, oy+sy)
		if err != nil {
			return err
		}
	}
	return nil
}

func Up(v View) error {
	if v == nil {
		return nil
	}
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	_, _, _, sy := v.Speed()
	err := v.SetCursor(cx, cy-sy)
	if err != nil && oy >= sy {
		err := v.SetOrigin(ox, oy-sy)
		if err != nil {
			return err
		}
	}
	return nil
}

func Right(v View) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	sx, _, _, _ := v.Speed()
	err := v.SetCursor(cx+sx, cy)
	if err != nil {
		ox, oy := v.Origin()
		err := v.SetOrigin(ox+sx, oy)
		if err != nil {
			return err
		}
	}
	return nil
}

func Left(v View) error {
	if v == nil {
		return nil
	}
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	_, sx, _, _ := v.Speed()
	err := v.SetCursor(cx-sx, cy)
	if err != nil {
		err := v.SetCursor(0, cy)
		if err != nil {
			return err
		}

		if ox >= sx-cx {
			err := v.SetOrigin(ox-sx+cx, oy)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
