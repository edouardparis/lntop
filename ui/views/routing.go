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
	ROUTING         = "routing"
	ROUTING_COLUMNS = "routing_columns"
	ROUTING_FOOTER  = "routing_footer"
)

var DefaultRoutingColumns = []string{
	"DIR",
	"STATUS",
	"IN_CHANNEL",
	"IN_ALIAS",
	"OUT_CHANNEL",
	"OUT_ALIAS",
	"AMOUNT",
	"FEE",
	"LAST UPDATE",
	"DETAIL",
}

type Routing struct {
	cfg *config.View

	columns []routingColumn

	columnHeadersView *gocui.View
	columnViews       []*gocui.View
	view              *gocui.View
	routingEvents     *models.RoutingLog

	ox, oy int
	cx, cy int
}

type routingColumn struct {
	name    string
	width   int
	display func(*netmodels.RoutingEvent, ...color.Option) string
}

func (c Routing) Name() string {
	return ROUTING
}

func (c *Routing) Wrap(v *gocui.View) View {
	c.view = v
	return c
}

func (c Routing) currentColumnIndex() int {
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

func (c Routing) Origin() (int, int) {
	return c.ox, c.oy
}

func (c Routing) Cursor() (int, int) {
	return c.cx, c.cy
}

func (c *Routing) SetCursor(cx, cy int) error {
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

func (c *Routing) SetOrigin(ox, oy int) error {
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

func (c *Routing) Speed() (int, int, int, int) {
	_, height := c.view.Size()
	current := c.currentColumnIndex()
	up := 0
	down := 0
	if c.Index() > 0 {
		up = 1
	}
	if c.Index() < len(c.routingEvents.Log)-1 && c.Index() < height {
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

func (c *Routing) Limits() (pageSize int, fullSize int) {
	_, pageSize = c.view.Size()
	fullSize = len(c.routingEvents.Log)
	if pageSize < fullSize {
		fullSize = pageSize
	}
	return
}

func (c Routing) Index() int {
	_, oy := c.Origin()
	_, cy := c.Cursor()
	return cy + oy
}

func (c *Routing) Delete(g *gocui.Gui) error {
	err := g.DeleteView(ROUTING_COLUMNS)
	if err != nil {
		return err
	}

	err = g.DeleteView(ROUTING)
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
	return g.DeleteView(ROUTING_FOOTER)
}

func (c *Routing) Set(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var err error
	setCursor := false
	c.columnHeadersView, err = g.SetView(ROUTING_COLUMNS, x0-1, y0, x1+2, y0+2, 0)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setCursor = true
	}
	c.columnHeadersView.Frame = false
	c.columnHeadersView.BgColor = gocui.ColorGreen
	c.columnHeadersView.FgColor = gocui.ColorBlack

	c.view, err = g.SetView(ROUTING, x0-1, y0+1, x1+2, y1-1, 0)
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
	c.view.Highlight = true
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

	footer, err := g.SetView(ROUTING_FOOTER, x0-1, y1-2, x1+2, y1, 0)
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
	fmt.Fprintln(footer, fmt.Sprintf("%s%s %s%s",
		blackBg("F2"), "Menu",
		blackBg("F10"), "Quit",
	))
	return nil
}

func (c *Routing) display(g *gocui.Gui) {
	c.columnHeadersView.Rewind()
	var buffer bytes.Buffer
	currentColumnIndex := c.currentColumnIndex()
	for i := range c.columns {
		if currentColumnIndex == i {
			buffer.WriteString(color.Cyan(color.Background)(c.columns[i].name))
			buffer.WriteString(" ")
			continue
		}
		buffer.WriteString(c.columns[i].name)
		buffer.WriteString(" ")
	}
	fmt.Fprintln(c.columnHeadersView, buffer.String())

	_, height := c.view.Size()
	numEvents := len(c.routingEvents.Log)

	j := 0
	if height < numEvents {
		j = numEvents - height
	}
	if len(c.columnViews) == 0 {
		c.columnViews = make([]*gocui.View, len(c.columns))
		x0, y0, _, y1 := c.view.Dimensions()
		for i := range c.columns {
			width := c.columns[i].width
			cc, _ := g.SetView("routing_content_"+c.columns[i].name, x0, y0, x0+width+2, y1, 0)
			cc.Frame = false
			cc.Autoscroll = false
			cc.SelBgColor = gocui.ColorCyan
			cc.SelFgColor = gocui.ColorBlack | gocui.AttrDim
			cc.Highlight = true
			c.columnViews[i] = cc
		}
	}
	rewind := true
	for ; j < numEvents; j++ {
		var item = c.routingEvents.Log[j]
		x0, y0, _, y1 := c.view.Dimensions()
		x0 -= c.ox
		for i := range c.columns {
			var opt color.Option
			if currentColumnIndex == i {
				opt = color.Bold
			}
			width := c.columns[i].width
			cc, _ := g.SetView("routing_content_"+c.columns[i].name, x0, y0, x0+width+2, y1, 0)
			c.columnViews[i] = cc
			if rewind {
				cc.Rewind()
			}
			fmt.Fprintln(cc, c.columns[i].display(item, opt), " ")
			x0 += width + 1
		}
		rewind = false
	}
}

func NewRouting(cfg *config.View, routingEvents *models.RoutingLog, channels *models.Channels) *Routing {
	routing := &Routing{
		cfg:           cfg,
		routingEvents: routingEvents,
	}

	printer := message.NewPrinter(language.English)

	columns := DefaultRoutingColumns
	if cfg != nil && len(cfg.Columns) != 0 {
		columns = cfg.Columns
	}

	routing.columns = make([]routingColumn, len(columns))

	for i := range columns {
		switch columns[i] {
		case "DIR":
			routing.columns[i] = routingColumn{
				width:   4,
				name:    fmt.Sprintf("%-4s", columns[i]),
				display: rdirection,
			}
		case "STATUS":
			routing.columns[i] = routingColumn{
				width:   8,
				name:    fmt.Sprintf("%-8s", columns[i]),
				display: rstatus,
			}
		case "IN_ALIAS":
			routing.columns[i] = routingColumn{
				width:   25,
				name:    fmt.Sprintf("%-25s", columns[i]),
				display: ralias(channels, false),
			}
		case "IN_CHANNEL":
			routing.columns[i] = routingColumn{
				width: 19,
				name:  fmt.Sprintf("%19s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					if c.IncomingChannelId == 0 {
						return fmt.Sprintf("%19s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%19d", c.IncomingChannelId))
				},
			}
		case "IN_SCID":
			routing.columns[i] = routingColumn{
				width: 14,
				name:  fmt.Sprintf("%14s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					if c.IncomingChannelId == 0 {
						return fmt.Sprintf("%14s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%14s", ToScid(c.IncomingChannelId)))
				},
			}
		case "IN_TIMELOCK":
			routing.columns[i] = routingColumn{
				width: 10,
				name:  fmt.Sprintf("%10s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					if c.IncomingTimelock == 0 {
						return fmt.Sprintf("%10s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%10d", c.IncomingTimelock))
				},
			}
		case "IN_HTLC":
			routing.columns[i] = routingColumn{
				width: 10,
				name:  fmt.Sprintf("%10s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					if c.IncomingHtlcId == 0 {
						return fmt.Sprintf("%10s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%10d", c.IncomingHtlcId))
				},
			}
		case "OUT_ALIAS":
			routing.columns[i] = routingColumn{
				width:   25,
				name:    fmt.Sprintf("%-25s", columns[i]),
				display: ralias(channels, true),
			}
		case "OUT_CHANNEL":
			routing.columns[i] = routingColumn{
				width: 19,
				name:  fmt.Sprintf("%19s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					if c.OutgoingChannelId == 0 {
						return fmt.Sprintf("%19s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%19d", c.OutgoingChannelId))
				},
			}
		case "OUT_SCID":
			routing.columns[i] = routingColumn{
				width: 14,
				name:  fmt.Sprintf("%14s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					if c.OutgoingChannelId == 0 {
						return fmt.Sprintf("%14s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%14s", ToScid(c.OutgoingChannelId)))
				},
			}
		case "OUT_TIMELOCK":
			routing.columns[i] = routingColumn{
				width: 10,
				name:  fmt.Sprintf("%10s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					if c.OutgoingTimelock == 0 {
						return fmt.Sprintf("%10s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%10d", c.OutgoingTimelock))
				},
			}
		case "OUT_HTLC":
			routing.columns[i] = routingColumn{
				width: 10,
				name:  fmt.Sprintf("%10s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					if c.OutgoingHtlcId == 0 {
						return fmt.Sprintf("%10s", "")
					}
					return color.White(opts...)(fmt.Sprintf("%10d", c.OutgoingHtlcId))
				},
			}
		case "AMOUNT":
			routing.columns[i] = routingColumn{
				width: 12,
				name:  fmt.Sprintf("%12s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					return color.Yellow(opts...)(printer.Sprintf("%12d", c.AmountMsat/1000))
				},
			}
		case "FEE":
			routing.columns[i] = routingColumn{
				width: 8,
				name:  fmt.Sprintf("%8s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					return color.Yellow(opts...)(printer.Sprintf("%8d", c.FeeMsat/1000))
				},
			}
		case "LAST UPDATE":
			routing.columns[i] = routingColumn{
				width: 15,
				name:  fmt.Sprintf("%-15s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					return color.Cyan(opts...)(
						fmt.Sprintf("%15s", c.LastUpdate.Format("15:04:05 Jan _2")),
					)
				},
			}
		case "DETAIL":
			routing.columns[i] = routingColumn{
				width: 80,
				name:  fmt.Sprintf("%-80s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					return color.Cyan(opts...)(fmt.Sprintf("%-80s", c.FailureDetail))
				},
			}
		default:
			routing.columns[i] = routingColumn{
				width: 10,
				name:  fmt.Sprintf("%-10s", columns[i]),
				display: func(c *netmodels.RoutingEvent, opts ...color.Option) string {
					return fmt.Sprintf("%-10s", "")
				},
			}
		}
	}

	return routing
}

func rstatus(c *netmodels.RoutingEvent, opts ...color.Option) string {
	switch c.Status {
	case netmodels.RoutingStatusActive:
		return color.Yellow(opts...)(fmt.Sprintf("%-8s", "active"))
	case netmodels.RoutingStatusSettled:
		return color.Green(opts...)(fmt.Sprintf("%-8s", "settled"))
	case netmodels.RoutingStatusFailed:
		return color.Red(opts...)(fmt.Sprintf("%-8s", "failed"))
	case netmodels.RoutingStatusLinkFailed:
		return color.Red(opts...)(fmt.Sprintf("%-8s", "linkfail"))
	}
	return ""
}

func rdirection(c *netmodels.RoutingEvent, opts ...color.Option) string {
	switch c.Direction {
	case netmodels.RoutingSend:
		return color.White(opts...)(fmt.Sprintf("%-4s", "send"))
	case netmodels.RoutingReceive:
		return color.White(opts...)(fmt.Sprintf("%-4s", "recv"))
	case netmodels.RoutingForward:
		return color.White(opts...)(fmt.Sprintf("%-4s", "forw"))
	}
	return "   "
}

func ralias(channels *models.Channels, out bool) func(*netmodels.RoutingEvent, ...color.Option) string {
	return func(c *netmodels.RoutingEvent, opts ...color.Option) string {
		id := c.IncomingChannelId
		if out {
			id = c.OutgoingChannelId
		}

		if id == 0 {
			return color.White(opts...)(fmt.Sprintf("%-25s", ""))
		}

		var alias string
		var forced bool
		aliasColor := color.White(opts...)
		for _, ch := range channels.List() {
			if ch.ID == id {
				alias, forced = ch.ShortAlias()
				if forced {
					aliasColor = color.Cyan(opts...)
				}
				break
			}
		}
		return aliasColor(fmt.Sprintf("%-25s", alias))
	}
}
