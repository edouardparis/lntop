package color

import "github.com/fatih/color"

type Color color.Color

var (
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	GreenBg = color.New(color.BgGreen, color.FgBlack).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	RedBg   = color.New(color.BgRed, color.FgBlack).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	CyanBg  = color.New(color.BgCyan, color.FgBlack).SprintFunc()
	WhiteBg = color.New(color.BgWhite, color.FgBlack).SprintFunc()
	BlackBg = color.New(color.BgBlack, color.FgWhite).SprintFunc()
)
