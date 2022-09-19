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

	if options.bg {
		return greenBg
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

func Magenta(opts ...Option) func(a ...interface{}) string {
	options := newOptions(opts)
	if options.bg {
		return magentaBg
	}
	return magentaBg
}

func HSL256(h, s, l float64, opts ...Option) func(a ...interface{}) string {
	options := newOptions(opts)
	val := color.HSL(h, s, l).C256().Value()
	c := color.S256(val)
	if options.bg {
		fg := color.White.C256().Value()
		if l > 0.5 {
			fg = color.Black.C256().Value()
		}
		c = color.S256(fg, val)
	}
	if options.bold {
		c.AddOpts(color.Bold)
	}
	return func(a ...interface{}) string {
		return c.Sprint(a...)
	}
}
