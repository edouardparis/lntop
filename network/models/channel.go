package models

import (
	"strings"
	"time"

	"github.com/edouardparis/lntop/logging"
	"github.com/mattn/go-runewidth"
)

const (
	ChannelActive = iota + 1
	ChannelInactive
	ChannelOpening
	ChannelClosing
	ChannelForceClosing
	ChannelWaitingClose
	ChannelClosed
)

type ChannelsBalance struct {
	Balance            int64
	PendingOpenBalance int64
}

func (m ChannelsBalance) MarshalLogObject(enc logging.ObjectEncoder) error {
	enc.AddInt64("balance", m.Balance)
	enc.AddInt64("pending_open_balance", m.PendingOpenBalance)

	return nil
}

type Channel struct {
	ID                  uint64
	Status              int
	RemotePubKey        string
	ChannelPoint        string
	Capacity            int64
	LocalBalance        int64
	RemoteBalance       int64
	CommitFee           int64
	CommitWeight        int64
	FeePerKiloWeight    int64
	UnsettledBalance    int64
	TotalAmountSent     int64
	TotalAmountReceived int64
	UpdatesCount        uint64
	CSVDelay            uint32
	Age                 uint32
	Private             bool
	PendingHTLC         []*HTLC
	LastUpdate          *time.Time
	Node                *Node
	LocalPolicy         *RoutingPolicy
	RemotePolicy        *RoutingPolicy
	BlocksTilMaturity   int32
}

func (m Channel) MarshalLogObject(enc logging.ObjectEncoder) error {
	enc.AddUint64("id", m.ID)
	enc.AddInt("status", m.Status)
	enc.AddString("remote_pubkey", m.RemotePubKey)
	enc.AddString("channel_point", m.ChannelPoint)
	enc.AddInt64("capacity", m.Capacity)
	enc.AddInt64("local_balance", m.LocalBalance)
	enc.AddInt64("remote_balance", m.RemoteBalance)
	enc.AddInt64("commit_fee", m.CommitFee)
	enc.AddInt64("commit_weight", m.CommitWeight)
	enc.AddInt64("unsettled_balance", m.UnsettledBalance)
	enc.AddInt64("total_amount_sent", m.TotalAmountSent)
	enc.AddInt64("total_amount_received", m.TotalAmountReceived)
	enc.AddUint64("updates_count", m.UpdatesCount)
	enc.AddBool("private", m.Private)

	return nil
}

func (m Channel) ShortAlias() (alias string, forced bool) {
	if m.Node != nil && m.Node.ForcedAlias != "" {
		alias = m.Node.ForcedAlias
		forced = true
	} else if m.Node == nil || m.Node.Alias == "" {
		alias = m.RemotePubKey[:25]
	} else {
		alias = strings.ReplaceAll(m.Node.Alias, "\ufe0f", "")
	}
	if runewidth.StringWidth(alias) > 25 {
		alias = runewidth.Truncate(alias, 25, "")
	}
	return
}

type ChannelUpdate struct {
}

type ChannelEdgeUpdate struct {
	ChanPoints []string
}

type RoutingPolicy struct {
	TimeLockDelta           uint32
	MinHtlc                 int64
	MaxHtlc                 uint64
	FeeBaseMsat             int64
	FeeRateMilliMsat        int64
	Disabled                bool
	InboundFeeBaseMsat      int32
	InboundFeeRateMilliMsat int32
}
