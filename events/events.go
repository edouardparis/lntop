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
)

type Event struct {
	Type string
	ID   string
}

func New(kind string) *Event {
	return &Event{Type: kind}
}
