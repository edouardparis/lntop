package models

import (
	"encoding/hex"

	"github.com/edouardparis/lntop/logging"
)

type ReceivedKind int

const (
	KindInvoice ReceivedKind = iota
	KindKeysend
)

type Invoice struct {
	// Index: index of this invoice.
	// Each newly created invoice will increment
	// this index making it monotonically increasing.
	Index  uint64
	Amount int64
	// AmountPaid: The amount that was accepted for this invoice, in satoshis.
	AmountPaid int64
	// AmountPaidInMSat: The amount that was accepted for this invoice, in milli satoshis.
	AmountPaidInMSat int64
	Description      string
	// RPreImage: The hex-encoded preimage (32 byte) which will allow
	// settling an incoming HTLC payable to this preimage
	RPreImage []byte
	// RHash: The hash of the preimage.
	RHash []byte
	// PaymentRequest: A bare-bones invoice for a payment within the Lightning Network.
	// With the details of the invoice, the sender has all the data
	// necessary to send a payment to the recipient.
	PaymentRequest  string
	DescriptionHash []byte
	// FallBackAddress: Fallback on-chain address.
	FallBackAddress string
	Settled         bool
	CreationDate    int64
	SettleDate      int64
	Expiry          int64
	// CLTVExpiry: Delta to use for the time-lock of the CLTV extended to the final hop.
	CLTVExpiry uint64
	// Private: Whether this invoice should include routing hints for private channels.
	Private bool
	// Kind indicates whether this was a regular invoice or a keysend (spontaneous) payment.
	Kind ReceivedKind
}

func (m Invoice) GetRHash() string {
	return hex.EncodeToString(m.RHash)
}

func (m Invoice) MarshalLogObject(enc logging.ObjectEncoder) error {
	enc.AddUint64("index", m.Index)
	enc.AddBool("private", m.Private)
	enc.AddInt64("amount", m.Amount)
	enc.AddInt64("amount_paid", m.AmountPaid)
	enc.AddString("r_hash", m.GetRHash())
	enc.AddString("description", m.Description)
	enc.AddString("r_pre_image", hex.EncodeToString(m.RPreImage))
	enc.AddString("payment_request", m.PaymentRequest)
	enc.AddBool("settled", m.Settled)
	enc.AddInt64("expiry", m.Expiry)

	return nil
}
