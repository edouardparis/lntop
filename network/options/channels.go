package options

type Channel func(*ChannelOptions)

type ChannelOptions struct {
	Active   bool
	Inactive bool
	Public   bool
	Private  bool
	Pending  bool
}

func WithChannelPending(c *ChannelOptions) { c.Pending = true }

func WithChannelPublic(v bool) Channel {
	return func(c *ChannelOptions) { c.Public = v }
}

func WithChannelPrivate(v bool) Channel {
	return func(c *ChannelOptions) { c.Private = v }
}

func WithChannelActive(v bool) Channel {
	return func(c *ChannelOptions) { c.Active = v }
}

func WithChannelInactive(v bool) Channel {
	return func(c *ChannelOptions) { c.Inactive = v }
}

func NewChannelOptions(options ...Channel) ChannelOptions {
	opts := ChannelOptions{}
	for i := range options {
		options[i](&opts)
	}
	return opts
}
