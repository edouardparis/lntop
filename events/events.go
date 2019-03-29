package events

const (
	Quit           = "quit"
	InvoiceCreated = "invoice.created"
	InvoiceSettled = "invoice.settled"
)

type Event struct {
	Type string
	ID   string
}

func New(kind string) *Event {
	return &Event{Type: kind}
}
