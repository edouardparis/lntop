package events

type Event string

const (
	BlockReceived         Event = "block.received"
	ChannelActive         Event = "channel.active"
	ChannelBalanceUpdated Event = "channel.balance.updated"
	ChannelInactive       Event = "channel.inactive"
	ChannelPending        Event = "channel.pending"
	InvoiceCreated        Event = "invoice.created"
	InvoiceSettled        Event = "invoice.settled"
	PeerUpdated           Event = "peer.updated"
	TransactionCreated    Event = "transaction.created"
	WalletBalanceUpdated  Event = "wallet.balance.updated"
)

type Publisher string

const (
	Channels     Publisher = "channels"
	Invoices     Publisher = "invoices"
	Transactions Publisher = "transactions"
)
