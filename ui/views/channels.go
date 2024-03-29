package views

import (
	"bytes"
	"fmt"

	"github.com/awesome-gocui/gocui"
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
	"SENT",
	"RECEIVED",
	"HTLC",
	"UNSETTLED",
	"CFEE",
	"LAST UPDATE",
	"PRIVATE",
	"ID",
}

type Channels struct {
	cfg *config.View

	columns []channelsColumn

	columnHeadersView *gocui.View
	columnViews       []*gocui.View
	view              *gocui.View

	channels *models.Channels

	ox, oy int
	cx, cy int
}

type channelsColumn struct {
	name    string
	width   int
	sorted  bool
	sort    func(models.Order) models.ChannelsSort
	display func(*netmodels.Channel, ...color.Option) string
}

func (c Channels) Name() string {
	return CHANNELS
}

func (c *Channels) Wrap(v *gocui.View) View {
	c.view = v
	return c
}

func (c Channels) currentColumnIndex() int {
	x := c.ox + c.cx
	index := 0
	sum := 0
	for i := range c.columns {
		sum += c.columns[i].width + 1
		if x < sum {
			return index
		}
		index++
	}
	return index
}

func (c Channels) Sort(column string, order models.Order) {
	if column == "" {
		index := c.currentColumnIndex()
		if index >= len(c.columns) {
			return
		}
		col := c.columns[index]
		if col.sort == nil {
			return
		}

		c.channels.Sort(col.sort(order))
		for i := range c.columns {
			c.columns[i].sorted = (i == index)
		}
	}
}

func (c Channels) Origin() (int, int) {
	return c.ox, c.oy
}

func (c Channels) Cursor() (int, int) {
	return c.cx, c.cy
}

func (c *Channels) SetCursor(cx, cy int) error {
	if err := cursorCompat(c.columnHeadersView, cx, 0); err != nil {
		return err
	}
	err := c.columnHeadersView.SetCursor(cx, 0)
	if err != nil {
		return err
	}

	for _, cv := range c.columnViews {
		if err := cursorCompat(c.view, cx, cy); err != nil {
			return err
		}
		err = cv.SetCursor(cx, cy)
		if err != nil {
			return err
		}
	}

	c.cx, c.cy = cx, cy
	return nil
}

func (c *Channels) SetOrigin(ox, oy int) error {
	err := c.columnHeadersView.SetOrigin(ox, 0)
	if err != nil {
		return err
	}
	err = c.view.SetOrigin(ox, oy)
	if err != nil {
		return err
	}

	for _, cv := range c.columnViews {
		err = cv.SetOrigin(0, oy)
		if err != nil {
			return err
		}
	}

	c.ox, c.oy = ox, oy
	return nil
}

func (c *Channels) Speed() (int, int, int, int) {
	current := c.currentColumnIndex()
	up := 0
	down := 0
	if c.Index() > 0 {
		up = 1
	}
	if c.Index() < c.channels.Len()-1 {
		down = 1
	}
	if current > len(c.columns)-1 {
		return 0, c.columns[current-1].width + 1, down, up
	}
	if current == 0 {
		return c.columns[0].width + 1, 0, down, up
	}
	return c.columns[current].width + 1,
		c.columns[current-1].width + 1,
		down, up
}

func (c *Channels) Limits() (pageSize int, fullSize int) {
	_, pageSize = c.view.Size()
	fullSize = c.channels.Len()
	return
}

func (c Channels) Index() int {
	_, oy := c.Origin()
	_, cy := c.Cursor()
	return cy + oy
}

func (c *Channels) Delete(g *gocui.Gui) error {
	err := g.DeleteView(CHANNELS_COLUMNS)
	if err != nil {
		return err
	}

	err = g.DeleteView(CHANNELS)
	if err != nil {
		return err
	}

	for _, cv := range c.columnViews {
		err = g.DeleteView(cv.Name())
		if err != nil {
			return err
		}
	}
	c.columnViews = c.columnViews[:0]
	return g.DeleteView(CHANNELS_FOOTER)
}

func (c *Channels) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	setCursor := false
	c.columnHeadersView, err = g.SetView(CHANNELS_COLUMNS, x0-1, y0, x1+2, y0+2, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setCursor = true
	}
	c.columnHeadersView.Frame = false
	c.columnHeadersView.BgColor = gocui.ColorGreen
	c.columnHeadersView.FgColor = gocui.ColorBlack

	c.view, err = g.SetView(CHANNELS, x0-1, y0+1, x1+2, y1-1, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setCursor = true
	}
	c.view.Frame = false
	c.view.Autoscroll = false
	c.view.SelBgColor = gocui.ColorCyan
	c.view.SelFgColor = gocui.ColorBlack | gocui.AttrDim
	c.view.Highlight = false
	c.display(g)

	if setCursor {
		ox, oy := c.Origin()
		err := c.SetOrigin(ox, oy)
		if err != nil {
			return err
		}

		cx, cy := c.Cursor()
		err = c.SetCursor(cx, cy)
		if err != nil {
			return err
		}
	}

	footer, err := g.SetView(CHANNELS_FOOTER, x0-1, y1-2, x1+2, y1, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	footer.Frame = false
	footer.BgColor = gocui.ColorCyan
	footer.FgColor = gocui.ColorBlack
	footer.Rewind()
	blackBg := color.Black(color.Background)
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s %s%s",
		blackBg("F2"), "Menu",
		blackBg("Enter"), "Channel",
		blackBg("F10"), "Quit",
	))
	return nil
}

func (c *Channels) display(g *gocui.Gui) {
	c.columnHeadersView.Rewind()
	var buffer bytes.Buffer
	currentColumnIndex := c.currentColumnIndex()
	for i := range c.columns {
		if currentColumnIndex == i {
			buffer.WriteString(color.Cyan(color.Background)(c.columns[i].name))
			buffer.WriteString(" ")
			continue
		} else if c.columns[i].sorted {
			buffer.WriteString(color.Magenta(color.Background)(c.columns[i].name))
			buffer.WriteString(" ")
			continue
		}
		buffer.WriteString(c.columns[i].name)
		buffer.WriteString(" ")
	}
	fmt.Fprintln(c.columnHeadersView, buffer.String())

	if len(c.columnViews) == 0 {
		c.columnViews = make([]*gocui.View, len(c.columns))
		x0, y0, _, y1 := c.view.Dimensions()
		for i := range c.columns {
			width := c.columns[i].width
			cc, _ := g.SetView("channel_content_"+c.columns[i].name, x0, y0, x0+width+2, y1, 0)
			cc.Frame = false
			cc.Autoscroll = false
			cc.SelBgColor = gocui.ColorCyan
			cc.SelFgColor = gocui.ColorBlack | gocui.AttrDim
			cc.Highlight = true
			c.columnViews[i] = cc
		}
	}
	for ci, item := range c.channels.List() {
		x0, y0, _, y1 := c.view.Dimensions()
		x0 -= c.ox
		for i := range c.columns {
			var opt color.Option
			if currentColumnIndex == i {
				opt = color.Bold
			}
			width := c.columns[i].width
			cc, _ := g.SetView("channel_content_"+c.columns[i].name, x0, y0, x0+width+2, y1, 0)
			c.columnViews[i] = cc
			if ci == 0 {
				cc.Rewind()
			}
			fmt.Fprintln(cc, c.columns[i].display(item, opt), " ")
			x0 += width + 1
		}
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
				width: 13,
				name:  fmt.Sprintf("%-13s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						// status meanings are kinda the opposite of their numerical value
						return models.IntSort(-c1.Status, -c2.Status, order)
					}
				},
				display: status,
			}
		case "ALIAS":
			channels.columns[i] = channelsColumn{
				width: 25,
				name:  fmt.Sprintf("%-25s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.StringSort(c1.Node.Alias, c2.Node.Alias, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					aliasColor := color.White(opts...)
					alias, forced := c.ShortAlias()
					if forced {
						aliasColor = color.Cyan(opts...)
					}
					return aliasColor(fmt.Sprintf("%-25s", alias))
				},
			}
		case "GAUGE":
			channels.columns[i] = channelsColumn{
				width: 21,
				name:  fmt.Sprintf("%-21s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.Float64Sort(
							float64(c1.LocalBalance)*100/float64(c1.Capacity),
							float64(c2.LocalBalance)*100/float64(c2.Capacity),
							order)
					}
				},
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
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.Int64Sort(c1.LocalBalance, c2.LocalBalance, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.Cyan(opts...)(printer.Sprintf("%12d", c.LocalBalance))
				},
			}
		case "REMOTE":
			channels.columns[i] = channelsColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.Int64Sort(c1.RemoteBalance, c2.RemoteBalance, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.Cyan(opts...)(printer.Sprintf("%12d", c.RemoteBalance))
				},
			}
		case "CAP":
			channels.columns[i] = channelsColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.Int64Sort(c1.Capacity, c2.Capacity, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.White(opts...)(printer.Sprintf("%12d", c.Capacity))
				},
			}
		case "SENT":
			channels.columns[i] = channelsColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.Int64Sort(c1.TotalAmountSent, c2.TotalAmountSent, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.Cyan(opts...)(printer.Sprintf("%12d", c.TotalAmountSent))
				},
			}
		case "RECEIVED":
			channels.columns[i] = channelsColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.Int64Sort(c1.TotalAmountReceived, c2.TotalAmountReceived, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.Cyan(opts...)(printer.Sprintf("%12d", c.TotalAmountReceived))
				},
			}
		case "HTLC":
			channels.columns[i] = channelsColumn{
				width: 5,
				name:  fmt.Sprintf("%5s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.IntSort(len(c1.PendingHTLC), len(c2.PendingHTLC), order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.Yellow(opts...)(fmt.Sprintf("%5d", len(c.PendingHTLC)))
				},
			}
		case "UNSETTLED":
			channels.columns[i] = channelsColumn{
				width: 10,
				name:  fmt.Sprintf("%-10s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.Int64Sort(c1.UnsettledBalance, c2.UnsettledBalance, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.Yellow(opts...)(printer.Sprintf("%10d", c.UnsettledBalance))
				},
			}
		case "CFEE":
			channels.columns[i] = channelsColumn{
				width: 6,
				name:  fmt.Sprintf("%-6s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.Int64Sort(c1.CommitFee, c2.CommitFee, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.White(opts...)(printer.Sprintf("%6d", c.CommitFee))
				},
			}
		case "LAST UPDATE":
			channels.columns[i] = channelsColumn{
				width: 15,
				name:  fmt.Sprintf("%-15s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.DateSort(c1.LastUpdate, c2.LastUpdate, order)
					}
				},
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
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						// public > private
						return models.BoolSort(!c1.Private, !c2.Private, order)
					}
				},
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
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.UInt64Sort(c1.ID, c2.ID, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					if c.ID == 0 {
						return fmt.Sprintf("%-19s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%-19d", c.ID))
				},
			}
		case "SCID":
			channels.columns[i] = channelsColumn{
				width: 14,
				name:  fmt.Sprintf("%-14s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.UInt64Sort(c1.ID, c2.ID, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					if c.ID == 0 {
						return fmt.Sprintf("%-14s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%-14s", ToScid(c.ID)))
				},
			}
		case "NUPD":
			channels.columns[i] = channelsColumn{
				width: 8,
				name:  fmt.Sprintf("%-8s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.UInt64Sort(c1.UpdatesCount, c2.UpdatesCount, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					return color.White(opts...)(printer.Sprintf("%8d", c.UpdatesCount))
				},
			}
		case "BASE_OUT":
			channels.columns[i] = channelsColumn{
				width: 8,
				name:  fmt.Sprintf("%-8s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						var c1f uint64
						var c2f uint64
						if c1.LocalPolicy != nil {
							c1f = uint64(c1.LocalPolicy.FeeBaseMsat)
						}
						if c2.LocalPolicy != nil {
							c2f = uint64(c2.LocalPolicy.FeeBaseMsat)
						}
						return models.UInt64Sort(c1f, c2f, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					var val int64
					if c.LocalPolicy != nil {
						val = c.LocalPolicy.FeeBaseMsat
					}
					return color.White(opts...)(printer.Sprintf("%8d", val))
				},
			}
		case "RATE_OUT":
			channels.columns[i] = channelsColumn{
				width: 8,
				name:  fmt.Sprintf("%-8s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						var c1f uint64
						var c2f uint64
						if c1.LocalPolicy != nil {
							c1f = uint64(c1.LocalPolicy.FeeRateMilliMsat)
						}
						if c2.LocalPolicy != nil {
							c2f = uint64(c2.LocalPolicy.FeeRateMilliMsat)
						}
						return models.UInt64Sort(c1f, c2f, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					var val int64
					if c.LocalPolicy != nil {
						val = c.LocalPolicy.FeeRateMilliMsat
					}
					return color.White(opts...)(printer.Sprintf("%8d", val))
				},
			}
		case "BASE_IN":
			channels.columns[i] = channelsColumn{
				width: 7,
				name:  fmt.Sprintf("%-7s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						var c1f uint64
						var c2f uint64
						if c1.RemotePolicy != nil {
							c1f = uint64(c1.RemotePolicy.FeeBaseMsat)
						}
						if c2.RemotePolicy != nil {
							c2f = uint64(c2.RemotePolicy.FeeBaseMsat)
						}
						return models.UInt64Sort(c1f, c2f, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					var val int64
					if c.RemotePolicy != nil {
						val = c.RemotePolicy.FeeBaseMsat
					}
					return color.White(opts...)(printer.Sprintf("%7d", val))
				},
			}
		case "RATE_IN":
			channels.columns[i] = channelsColumn{
				width: 7,
				name:  fmt.Sprintf("%-7s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						var c1f uint64
						var c2f uint64
						if c1.RemotePolicy != nil {
							c1f = uint64(c1.RemotePolicy.FeeRateMilliMsat)
						}
						if c2.RemotePolicy != nil {
							c2f = uint64(c2.RemotePolicy.FeeRateMilliMsat)
						}
						return models.UInt64Sort(c1f, c2f, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					var val int64
					if c.RemotePolicy != nil {
						val = c.RemotePolicy.FeeRateMilliMsat
					}
					return color.White(opts...)(printer.Sprintf("%7d", val))
				},
			}
		case "AGE":
			channels.columns[i] = channelsColumn{
				width: 10,
				name:  fmt.Sprintf("%10s", columns[i]),
				sort: func(order models.Order) models.ChannelsSort {
					return func(c1, c2 *netmodels.Channel) bool {
						return models.UInt32Sort(c1.Age, c2.Age, order)
					}
				},
				display: func(c *netmodels.Channel, opts ...color.Option) string {
					if c.ID == 0 {
						return fmt.Sprintf("%10s", "")
					}
					result := printer.Sprintf("%10s", FormatAge(c.Age))
					if cfg.Options.GetOption("AGE", "color") == "color" {
						return ColorizeAge(c.Age, result, opts...)
					} else {
						return color.White(opts...)(result)
					}
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

func channelDisabled(c *netmodels.Channel, opts ...color.Option) string {
	outgoing := false
	incoming := false
	if c.LocalPolicy != nil && c.LocalPolicy.Disabled {
		outgoing = true
	}
	if c.RemotePolicy != nil && c.RemotePolicy.Disabled {
		incoming = true
	}
	result := ""
	if incoming && outgoing {
		result = "⇅"
	} else if incoming {
		result = "⇊"
	} else if outgoing {
		result = "⇈"
	}
	if result == "" {
		return result
	}
	return color.Red(opts...)(fmt.Sprintf("%-4s", result))
}

func status(c *netmodels.Channel, opts ...color.Option) string {
	disabled := channelDisabled(c, opts...)
	format := "%-13s"
	if disabled != "" {
		format = "%-9s"
	}
	switch c.Status {
	case netmodels.ChannelActive:
		return color.Green(opts...)(fmt.Sprintf(format, "active ")) + disabled
	case netmodels.ChannelInactive:
		return color.Red(opts...)(fmt.Sprintf(format, "inactive ")) + disabled
	case netmodels.ChannelOpening:
		return color.Yellow(opts...)(fmt.Sprintf("%-13s", "opening"))
	case netmodels.ChannelClosing:
		return color.Yellow(opts...)(fmt.Sprintf("%-13s", "closing"))
	case netmodels.ChannelForceClosing:
		return color.Yellow(opts...)(fmt.Sprintf("%-13s", "force closing"))
	case netmodels.ChannelWaitingClose:
		return color.Yellow(opts...)(fmt.Sprintf("%-13s", "waiting close"))
	case netmodels.ChannelClosed:
		return color.Red(opts...)(fmt.Sprintf("%-13s", "closed"))
	}
	return ""
}
