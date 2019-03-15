package models

import "github.com/edouardparis/lntop/logging"

type WalletBalance struct {
	TotalBalance       int64
	ConfirmedBalance   int64
	UnconfirmedBalance int64
}

func (m WalletBalance) MarshalLogObject(enc logging.ObjectEncoder) error {
	enc.AddInt64("total_balance", m.TotalBalance)
	enc.AddInt64("confirmed_balance", m.ConfirmedBalance)
	enc.AddInt64("unconfirmed_balance", m.UnconfirmedBalance)

	return nil
}
