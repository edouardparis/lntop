package color

import "github.com/gookit/color"

type Color color.Color

var (
	yellow     = SprintFunc(color.New(color.FgYellow))
	yellowBold = SprintFunc(color.New(color.FgYellow, color.Bold))
	green      = SprintFunc(color.New(color.FgGreen))
	greenBold  = SprintFunc(color.New(color.FgGreen, color.Bold))
	greenBg    = SprintFunc(color.New(color.FgBlack, color.BgGreen))
	magentaBg  = SprintFunc(color.New(color.FgBlack, color.BgMagenta))
	red        = SprintFunc(color.New(color.FgRed))
	redBold    = SprintFunc(color.New(color.FgRed, color.Bold))
	cyan       = SprintFunc(color.New(color.FgCyan))
	cyanBold   = SprintFunc(color.New(color.FgCyan, color.Bold))
	cyanBg     = SprintFunc(color.New(color.BgCyan, color.FgBlack))
	white      = SprintFunc(color.New())
	whiteBold  = SprintFunc(color.New(color.Bold))
	blackBg    = SprintFunc(color.New(color.BgBlack, color.FgWhite))
	black      = SprintFunc(color.New(color.FgBlack))
)

func SprintFunc(c color.Style) func(args ...interface{}) string {
	return func(args ...interface{}) string {
		return c.Sprint(args...)
	}
}

type Option func(*options)

type options struct {
	Bold bool
	Bg   bool
}

func NewOptions(opts []Option) options {
	options := options{}
	for i := range opts {
		if opts[i] == nil {
			continue
		}
		opts[i](&options)
	}
	return options
}

func Bold(o *options)       { o.Bold = true }
func Background(o *options) { o.Bg = true }

func Yellow(opts ...Option) func(a ...interface{}) string {
	options := NewOptions(opts)
	if options.Bold {
		return yellowBold
	}
	return yellow
}

func Green(opts ...Option) func(a ...interface{}) string {
	options := NewOptions(opts)
	if options.Bold {
		return greenBold
	}

	if options.Bg {
		return greenBg
	}

	return green
}

func Red(opts ...Option) func(a ...interface{}) string {
	options := NewOptions(opts)
	if options.Bold {
		return redBold
	}
	return red
}

func White(opts ...Option) func(a ...interface{}) string {
	options := NewOptions(opts)
	if options.Bold {
		return whiteBold
	}
	return white
}

func Cyan(opts ...Option) func(a ...interface{}) string {
	options := NewOptions(opts)
	if options.Bold {
		return cyanBold
	}
	if options.Bg {
		return cyanBg
	}
	return cyan
}

func Black(opts ...Option) func(a ...interface{}) string {
	options := NewOptions(opts)
	if options.Bg {
		return blackBg
	}
	return black
}

func Magenta(opts ...Option) func(a ...interface{}) string {
	options := NewOptions(opts)
	if options.Bg {
		return magentaBg
	}
	return magentaBg
}
