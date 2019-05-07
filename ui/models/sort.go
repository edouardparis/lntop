package models

type Order int

const (
	Asc Order = iota
	Desc
)

func Int64Sort(a, b int64, o Order) bool {
	if o == Asc {
		return a > b
	}
	return a < b
}
