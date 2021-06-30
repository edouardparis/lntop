package cursor

type View interface {
	Cursor() (int, int)
	Origin() (int, int)
	Speed() (right int, left int, down int, up int)
	Limits() (pageSize int, fullSize int)
	SetCursor(int, int) error
	SetOrigin(int, int) error
}

func Down(v View) error {
	if v == nil {
		return nil
	}
	cx, cy := v.Cursor()
	ox, oy := v.Origin()
	_, _, sy, _ := v.Speed()
	_, fs := v.Limits()
	if cy+oy+sy >= fs {
		return nil
	}
	err := v.SetCursor(cx, cy+sy)
	if err != nil {
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

func Home(v View) error {
	if v == nil {
		return nil
	}
	ox, _ := v.Origin()
	cx, _ := v.Cursor()
	v.SetCursor(cx, 0)
	v.SetOrigin(ox, 0)
	return nil
}

func End(v View) error {
	if v == nil {
		return nil
	}
	ps, fs := v.Limits()
	if ps == 0 { // no pagination
		return nil
	}
	if ps > fs {
		ps = fs
	}
	ox, _ := v.Origin()
	cx, _ := v.Cursor()
	v.SetCursor(cx, ps-1)
	v.SetOrigin(ox, fs-ps)
	return nil
}

func PageDown(v View) error {
	if v == nil {
		return nil
	}
	ps, fs := v.Limits()
	if ps == 0 { // no pagination
		return nil
	}
	if ps > fs {
		ps = fs
	}
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	ny := oy + cy + ps
	if ny >= fs {
		ny = fs - 1
	}
	if ny >= fs-ps {
		v.SetOrigin(ox, fs-ps)
		v.SetCursor(cx, ny-fs+ps)
	} else {
		v.SetOrigin(ox, ny-ps)
		v.SetCursor(cx, ps-1)
	}
	return nil
}

func PageUp(v View) error {
	if v == nil {
		return nil
	}
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	ps, _ := v.Limits()
	ny := oy + cy - ps
	if ny <= 0 {
		ny = 0
	}
	if ny <= ps {
		v.SetOrigin(ox, 0)
		v.SetCursor(cx, ny)
	} else {
		v.SetOrigin(ox, ny)
		v.SetCursor(cx, 0)
	}
	return nil
}
