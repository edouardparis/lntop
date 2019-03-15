package models

import "github.com/edouardparis/lntop/logging"

type Payment struct {
	PaymentError    string
	PaymentPreimage []byte
	PayReq          *PayReq
	Route           *Route
}

func (p Payment) MarshalLogObject(enc logging.ObjectEncoder) error {
	enc.AddString("payment_error", p.PaymentError)

	return nil
}
