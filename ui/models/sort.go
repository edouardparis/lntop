package models

import "time"

type Order int

const (
	Asc Order = iota
	Desc
)

func IntSort(a, b int, o Order) bool {
	if o == Asc {
		return a < b
	}
	return a > b
}

func Int32Sort(a, b int32, o Order) bool {
	if o == Asc {
		return a < b
	}
	return a > b
}

func Int64Sort(a, b int64, o Order) bool {
	if o == Asc {
		return a < b
	}
	return a > b
}

func DateSort(a, b *time.Time, o Order) bool {
	if o == Desc {
		if a == nil || b == nil {
			return b == nil
		}

		return a.After(*b)
	}

	if a == nil || b == nil {
		return a == nil
	}

	return a.Before(*b)
}
