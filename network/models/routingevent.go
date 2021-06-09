package models

import (
	"time"
)

const (
	RoutingSend = iota + 1
	RoutingReceive
	RoutingForward
)

const (
	RoutingStatusActive = iota + 1
	RoutingStatusFailed
	RoutingStatusSettled
	RoutingStatusLinkFailed
)

type RoutingEvent struct {
	IncomingChannelId uint64
	OutgoingChannelId uint64
	IncomingHtlcId    uint64
	OutgoingHtlcId    uint64
	LastUpdate        time.Time
	Direction         int
	Status            int
	IncomingTimelock  uint32
	OutgoingTimelock  uint32
	AmountMsat        uint64
	FeeMsat           uint64
	FailureCode       int32
	FailureDetail     string
}

func (u *RoutingEvent) Equals(other *RoutingEvent) bool {
	return u.IncomingChannelId == other.IncomingChannelId && u.IncomingHtlcId == other.IncomingHtlcId && u.OutgoingChannelId == other.OutgoingChannelId && u.OutgoingHtlcId == other.OutgoingHtlcId
}

func (u *RoutingEvent) Update(newer *RoutingEvent) {
	u.LastUpdate = newer.LastUpdate
	u.Status = newer.Status
	u.FailureCode = newer.FailureCode
	u.FailureDetail = newer.FailureDetail
}
