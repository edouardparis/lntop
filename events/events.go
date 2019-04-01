package events

const (
	PeerUpdated     = "peer.updated"
	BlockReceived   = "block.received"
	InvoiceCreated  = "invoice.created"
	InvoiceSettled  = "invoice.settled"
	ChannelPending  = "channel.pending"
	ChannelActive   = "channel.active"
	ChannelInactive = "channel.inactive"
)

type Event struct {
	Type string
	ID   string
}

func New(kind string) *Event {
	return &Event{Type: kind}
}
