package color

import "github.com/fatih/color"

type Color color.Color

var (
	yellow     = color.New(color.FgYellow).SprintFunc()
	yellowBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	green      = color.New(color.FgGreen).SprintFunc()
	greenBold  = color.New(color.FgGreen, color.Bold).SprintFunc()
	red        = color.New(color.FgRed).SprintFunc()
	redBold    = color.New(color.FgRed, color.Bold).SprintFunc()
	cyan       = color.New(color.FgCyan).SprintFunc()
	cyanBold   = color.New(color.FgCyan, color.Bold).SprintFunc()
	cyanBg     = color.New(color.BgCyan, color.FgBlack).SprintFunc()
	white      = color.New(color.FgWhite).SprintFunc()
	whiteBold  = color.New(color.FgWhite, color.Bold).SprintFunc()
	blackBg    = color.New(color.BgBlack, color.FgWhite).SprintFunc()
	black      = color.New(color.FgBlack).SprintFunc()
)

type Option func(*options)

type options struct {
	bold bool
	bg   bool
}

func newOptions(opts []Option) options {
	options := options{}
	for i := range opts {
		if opts[i] == nil {
			continue
		}
		opts[i](&options)
	}
	return options
}

func Bold(o *options)       { o.bold = true }
func Background(o *options) { o.bg = true }

func Yellow(opts ...Option) func(a ...interface{}) string {
	options := newOptions(opts)
	if options.bold {
		return yellowBold
	}
	return yellow
}

func Green(opts ...Option) func(a ...interface{}) string {
	options := newOptions(opts)
	if options.bold {
		return greenBold
	}
	return green
}

func Red(opts ...Option) func(a ...interface{}) string {
	options := newOptions(opts)
	if options.bold {
		return redBold
	}
	return red
}

func White(opts ...Option) func(a ...interface{}) string {
	options := newOptions(opts)
	if options.bold {
		return whiteBold
	}
	return white
}

func Cyan(opts ...Option) func(a ...interface{}) string {
	options := newOptions(opts)
	if options.bold {
		return cyanBold
	}
	if options.bg {
		return cyanBg
	}
	return cyan
}

func Black(opts ...Option) func(a ...interface{}) string {
	options := newOptions(opts)
	if options.bg {
		return blackBg
	}
	return black
}
