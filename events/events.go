package events

const (
	BlockReceived         = "block.received"
	ChannelActive         = "channel.active"
	ChannelBalanceUpdated = "channel.balance.updated"
	ChannelInactive       = "channel.inactive"
	ChannelPending        = "channel.pending"
	InvoiceCreated        = "invoice.created"
	InvoiceSettled        = "invoice.settled"
	PeerUpdated           = "peer.updated"
	TransactionCreated    = "transaction.created"
	WalletBalanceUpdated  = "wallet.balance.updated"
	RoutingEventUpdated   = "routing.event.updated"
	GraphUpdated          = "graph.updated"
)

type Event struct {
	Type string
	ID   string
	Data interface{}
}

func New(kind string) *Event {
	return &Event{Type: kind}
}

func NewWithData(kind string, data interface{}) *Event {
	return &Event{Type: kind, Data: data}
}
