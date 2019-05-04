package views

import (
	"bytes"
	"fmt"

	"github.com/jroimartin/gocui"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/edouardparis/lntop/config"
	netmodels "github.com/edouardparis/lntop/network/models"
	"github.com/edouardparis/lntop/ui/color"
	"github.com/edouardparis/lntop/ui/models"
)

const (
	CHANNELS         = "channels"
	CHANNELS_COLUMNS = "channels_columns"
	CHANNELS_FOOTER  = "channels_footer"
)

var DefaultChannelsColumns = []string{
	"STATUS",
	"ALIAS",
	"GAUGE",
	"LOCAL",
	"CAP",
	"HTLC",
	"UNSETTLED",
	"CFEE",
	"LAST UPDATE",
	"PRIVATE",
	"ID",
}

type Channels struct {
	cfg *config.View

	col     int
	columns []channelsColumn

	columnsView *gocui.View
	view        *gocui.View
	channels    *models.Channels
}

type channelsColumn struct {
	name    string
	width   int
	display func(*netmodels.Channel, ...color.Option) string
}

func (c Channels) Index() int {
	_, oy := c.view.Origin()
	_, cy := c.view.Cursor()
	return cy + oy
}

func (c Channels) Name() string {
	return CHANNELS
}

func (c *Channels) Wrap(v *gocui.View) view {
	c.view = v
	return c
}

func (c *Channels) CursorDown() error {
	return cursorDown(c.view, 1)
}

func (c *Channels) CursorUp() error {
	return cursorUp(c.view, 1)
}

func (c *Channels) CursorRight() error {
	if c.col >= len(c.columns)-1 {
		return nil
	}
	speed := c.columns[c.col].width + 1
	c.col++
	err := cursorRight(c.columnsView, speed)
	if err != nil {
		return err
	}

	return cursorRight(c.view, speed)
}

func (c *Channels) CursorLeft() error {
	if c.col == 0 {
		return nil
	}
	speed := c.columns[c.col-1].width + 1
	c.col--
	err := cursorLeft(c.columnsView, speed)
	if err != nil {
		return err
	}

	return cursorLeft(c.view, speed)
}

func (c Channels) Delete(g *gocui.Gui) error {
	err := g.DeleteView(CHANNELS_COLUMNS)
	if err != nil {
		return err
	}

	err = g.DeleteView(CHANNELS)
	if err != nil {
		return err
	}

	return g.DeleteView(CHANNELS_FOOTER)
}

func (c *Channels) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	c.columnsView, err = g.SetView(CHANNELS_COLUMNS, x0-1, y0, x1+2, y0+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	c.columnsView.Frame = false
	c.columnsView.BgColor = gocui.ColorGreen
	c.columnsView.FgColor = gocui.ColorBlack

	c.view, err = g.SetView(CHANNELS, x0-1, y0+1, x1+2, y1-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	c.view.Frame = false
	c.view.Autoscroll = false
	c.view.SelBgColor = gocui.ColorCyan
	c.view.SelFgColor = gocui.ColorBlack
	c.view.Highlight = true

	c.display()

	footer, err := g.SetView(CHANNELS_FOOTER, x0-1, y1-2, x1+2, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	footer.Frame = false
	footer.BgColor = gocui.ColorCyan
	footer.FgColor = gocui.ColorBlack
	footer.Clear()
	blackBg := color.Black(color.Background)
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s %s%s %s%s",
		blackBg("F1"), "Help",
		blackBg("F2"), "Menu",
		blackBg("Enter"), "Channel",
		blackBg("F10"), "Quit",
	))
	return nil
}

func (c *Channels) display() {
	c.columnsView.Clear()
	var buffer bytes.Buffer
	for i := range c.columns {
		if c.col == i {
			buffer.WriteString(color.White(color.Bold)(c.columns[i].name))
			buffer.WriteString(" ")
			continue
		}
		buffer.WriteString(c.columns[i].name)
		buffer.WriteString(" ")
	}
	fmt.Fprintln(c.columnsView, buffer.String())

	c.view.Clear()
	for _, item := range c.channels.List() {
		var buffer bytes.Buffer
		for i := range c.columns {
			var opt color.Option
			if c.col == i {
				opt = color.Bold
			}
			buffer.WriteString(c.columns[i].display(item, opt))
			buffer.WriteString(" ")
		}
		fmt.Fprintln(c.view, buffer.String())
	}
}

func NewChannels(cfg *config.View, chans *models.Channels) *Channels {
	channels := &Channels{
		cfg:      cfg,
		channels: chans,
	}

	printer := message.NewPrinter(language.English)

	columns := DefaultChannelsColumns
	if cfg != nil && len(cfg.Columns) != 0 {
		columns = cfg.Columns
	}

	channels.columns = make([]channelsColumn, len(columns))

	for i := range columns {
		switch columns[i] {
		case "STATUS":
			channels.columns[i] = channelsColumn{
				width:   13,
				name:    fmt.Sprintf("%-13s", columns[i]),
				display: status,
			}
		case "ALIAS":
			channels.columns[i] = channelsColumn{
				width: 25,
				name:  fmt.Sprintf("%-25s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					var alias string
					if c.Node == nil || c.Node.Alias == "" {
						alias = c.RemotePubKey[:24]
					} else if len(c.Node.Alias) > 25 {
						alias = c.Node.Alias[:24]
					} else {
						alias = c.Node.Alias
					}
					return color.White(opts...)(fmt.Sprintf("%-25s", alias))
				},
			}
		case "GAUGE":
			channels.columns[i] = channelsColumn{
				width: 21,
				name:  fmt.Sprintf("%-21s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					index := int(c.LocalBalance * int64(15) / c.Capacity)
					var buffer bytes.Buffer
					cyan := color.Cyan(opts...)
					white := color.White(opts...)
					for i := 0; i < 15; i++ {
						if i < index {
							buffer.WriteString(cyan("|"))
							continue
						}
						buffer.WriteString(" ")
					}
					return fmt.Sprintf("%s%s%s",
						white("["),
						buffer.String(),
						white(fmt.Sprintf("] %2d%%", c.LocalBalance*100/c.Capacity)))
				},
			}
		case "LOCAL":
			channels.columns[i] = channelsColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.Cyan(opts...)(printer.Sprintf("%12d", c.LocalBalance))
				},
			}
		case "CAP":
			channels.columns[i] = channelsColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.White(opts...)(printer.Sprintf("%12d", c.Capacity))
				},
			}
		case "HTLC":
			channels.columns[i] = channelsColumn{
				width: 5,
				name:  fmt.Sprintf("%5s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.Yellow(opts...)(fmt.Sprintf("%5d", len(c.PendingHTLC)))
				},
			}
		case "UNSETTLED":
			channels.columns[i] = channelsColumn{
				width: 10,
				name:  fmt.Sprintf("%-10s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.Yellow(opts...)(printer.Sprintf("%10d", c.UnsettledBalance))
				},
			}
		case "CFEE":
			channels.columns[i] = channelsColumn{
				width: 6,
				name:  fmt.Sprintf("%-6s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.White(opts...)(printer.Sprintf("%6d", c.CommitFee))
				},
			}
		case "LAST UPDATE":
			channels.columns[i] = channelsColumn{
				width: 15,
				name:  fmt.Sprintf("%-15s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					if c.LastUpdate != nil {
						return color.Cyan(opts...)(
							fmt.Sprintf("%15s", c.LastUpdate.Format("15:04:05 Jan _2")),
						)
					}
					return fmt.Sprintf("%15s", "")
				},
			}
		case "PRIVATE":
			channels.columns[i] = channelsColumn{
				width: 7,
				name:  fmt.Sprintf("%-7s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					if c.Private {
						return color.Red(opts...)("private")
					}
					return color.Green(opts...)("public ")
				},
			}
		case "ID":
			channels.columns[i] = channelsColumn{
				width: 19,
				name:  fmt.Sprintf("%-19s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					if c.ID == 0 {
						return fmt.Sprintf("%-19s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%d", c.ID))
				},
			}
		default:
			channels.columns[i] = channelsColumn{
				width: 21,
				name:  fmt.Sprintf("%-21s", columns[i]),
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return "column does not exist"
				},
			}
		}

	}
	return channels
}

func status(c *netmodels.Channel, opts ...color.Option) string {
	switch c.Status {
	case netmodels.ChannelActive:
		return color.Green(opts...)(fmt.Sprintf("%-13s", "active"))
	case netmodels.ChannelInactive:
		return color.Red(opts...)(fmt.Sprintf("%-13s", "inactive"))
	case netmodels.ChannelOpening:
		return color.Yellow(opts...)(fmt.Sprintf("%-13s", "opening"))
	case netmodels.ChannelClosing:
		return color.Yellow(opts...)(fmt.Sprintf("%-13s", "closing"))
	case netmodels.ChannelForceClosing:
		return color.Yellow(opts...)(fmt.Sprintf("%-13s", "force closing"))
	case netmodels.ChannelWaitingClose:
		return color.Yellow(opts...)(fmt.Sprintf("%-13s", "waiting close"))
	}
	return ""
}
