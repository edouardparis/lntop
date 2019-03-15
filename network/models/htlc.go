package models

type HTLC struct {
	Incoming         bool
	Amount           int64
	Hashlock         []byte
	ExpirationHeight uint32
}
