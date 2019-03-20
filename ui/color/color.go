package color

import "github.com/fatih/color"

type Color color.Color

var (
	Green   = color.New(color.FgGreen).SprintFunc()
	GreenBg = color.New(color.BgGreen, color.FgBlack).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	CyanBg  = color.New(color.BgCyan, color.FgBlack).SprintFunc()
)
