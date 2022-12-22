package models

import "time"

type ForwardingEvent struct {
	PeerAliasIn  string
	PeerAliasOut string
	ChanIdIn     uint64
	ChanIdOut    uint64
	AmtIn        uint64
	AmtOut       uint64
	Fee          uint64
	FeeMsat      uint64
	AmtInMsat    uint64
	AmtOutMsat   uint64
	EventTime    time.Time
}
