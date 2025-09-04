package lnd

import (
	"fmt"
	"strings"
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/routerrpc"

	"github.com/edouardparis/lntop/network/models"
)

func protoToWalletBalance(w *lnrpc.WalletBalanceResponse) *models.WalletBalance {
	return &models.WalletBalance{
		TotalBalance:       w.GetTotalBalance(),
		ConfirmedBalance:   w.GetConfirmedBalance(),
		UnconfirmedBalance: w.GetUnconfirmedBalance(),
	}
}

func protoToChannelsBalance(w *lnrpc.ChannelBalanceResponse) *models.ChannelsBalance {
	return &models.ChannelsBalance{
		PendingOpenBalance: w.GetPendingOpenBalance(),
		Balance:            w.GetBalance(),
	}
}

func addInvoiceProtoToInvoice(req *lnrpc.Invoice, resp *lnrpc.AddInvoiceResponse) *models.Invoice {
	inv := &models.Invoice{
		Expiry:         req.GetExpiry(),
		Amount:         req.GetValue(),
		Description:    req.GetMemo(),
		CreationDate:   req.GetCreationDate(),
		RHash:          resp.GetRHash(),
		PaymentRequest: resp.GetPaymentRequest(),
		Index:          resp.GetAddIndex(),
	}
	// Determine kind: if PaymentRequest is empty, treat as keysend
	if inv.PaymentRequest == "" {
		inv.Kind = models.KindKeysend
	} else {
		inv.Kind = models.KindInvoice
	}
	return inv
}

func lookupInvoiceProtoToInvoice(resp *lnrpc.Invoice) *models.Invoice {
	inv := &models.Invoice{
		Index:            resp.GetAddIndex(),
		Amount:           resp.GetValue(),
		AmountPaid:       resp.GetAmtPaidSat(),
		AmountPaidInMSat: resp.GetAmtPaidMsat(),
		Description:      resp.GetMemo(),
		RPreImage:        resp.GetRPreimage(),
		RHash:            resp.GetRHash(),
		PaymentRequest:   resp.GetPaymentRequest(),
		DescriptionHash:  resp.GetDescriptionHash(),
		FallBackAddress:  resp.GetFallbackAddr(),
		Settled:          resp.GetSettled(),
		CreationDate:     resp.GetCreationDate(),
		SettleDate:       resp.GetSettleDate(),
		Expiry:           resp.GetExpiry(),
		CLTVExpiry:       resp.GetCltvExpiry(),
		Private:          resp.GetPrivate(),
	}
	if inv.PaymentRequest == "" {
		inv.Kind = models.KindKeysend
	} else {
		inv.Kind = models.KindInvoice
	}
	return inv
}

func listChannelsProtoToChannels(r *lnrpc.ListChannelsResponse) []*models.Channel {
	resp := r.GetChannels()
	channels := make([]*models.Channel, len(resp))
	for i := range resp {
		channels[i] = channelProtoToChannel(resp[i])
	}

	return channels
}

func channelProtoToChannel(c *lnrpc.Channel) *models.Channel {
	htlcs := c.GetPendingHtlcs()
	HTLCs := make([]*models.HTLC, len(htlcs))
	for i := range htlcs {
		HTLCs[i] = htlcProtoToHTLC(htlcs[i])
	}

	status := models.ChannelInactive
	if c.Active {
		status = models.ChannelActive
	}

	return &models.Channel{
		ID:                  c.GetChanId(),
		Status:              status,
		RemotePubKey:        c.GetRemotePubkey(),
		ChannelPoint:        c.GetChannelPoint(),
		Capacity:            c.GetCapacity(),
		LocalBalance:        c.GetLocalBalance(),
		RemoteBalance:       c.GetRemoteBalance(),
		CommitFee:           c.GetCommitFee(),
		CommitWeight:        c.GetCommitWeight(),
		FeePerKiloWeight:    c.GetFeePerKw(),
		UnsettledBalance:    c.GetUnsettledBalance(),
		TotalAmountSent:     c.GetTotalSatoshisSent(),
		TotalAmountReceived: c.GetTotalSatoshisReceived(),
		UpdatesCount:        c.GetNumUpdates(),
		CSVDelay:            c.GetCsvDelay(),
		Private:             c.GetPrivate(),
		PendingHTLC:         HTLCs,
	}
}

func htlcProtoToHTLC(h *lnrpc.HTLC) *models.HTLC {
	return &models.HTLC{
		Incoming:         h.GetIncoming(),
		Amount:           h.GetAmount(),
		Hashlock:         h.GetHashLock(),
		ExpirationHeight: h.GetExpirationHeight(),
	}
}

func pendingChannelsProtoToChannels(r *lnrpc.PendingChannelsResponse) []*models.Channel {
	respPending := r.GetPendingOpenChannels()
	pending := make([]*models.Channel, len(respPending))
	for i := range respPending {
		pending[i] = openingChannelProtoToChannel(respPending[i])
	}

	respClosing := r.GetPendingClosingChannels()
	closing := make([]*models.Channel, len(respClosing))
	for i := range respClosing {
		closing[i] = closingChannelProtoToChannel(respClosing[i])
	}

	channels := append(pending, closing...)

	respForceClosing := r.GetPendingForceClosingChannels()
	forceClosing := make([]*models.Channel, len(respForceClosing))
	for i := range respForceClosing {
		forceClosing[i] = forceClosingChannelProtoToChannel(respForceClosing[i])
	}

	channels = append(channels, forceClosing...)

	respWaitingClose := r.GetWaitingCloseChannels()
	waitingClose := make([]*models.Channel, len(respWaitingClose))
	for i := range respWaitingClose {
		waitingClose[i] = waitingCloseChannelProtoToChannel(respWaitingClose[i])
	}

	return append(channels, waitingClose...)
}

func openingChannelProtoToChannel(c *lnrpc.PendingChannelsResponse_PendingOpenChannel) *models.Channel {
	return &models.Channel{
		Status:           models.ChannelOpening,
		RemotePubKey:     c.Channel.RemoteNodePub,
		Capacity:         c.Channel.Capacity,
		LocalBalance:     c.Channel.LocalBalance,
		RemoteBalance:    c.Channel.RemoteBalance,
		ChannelPoint:     c.Channel.ChannelPoint,
		CommitWeight:     c.CommitWeight,
		CommitFee:        c.CommitFee,
		FeePerKiloWeight: c.FeePerKw,
	}
}

func closingChannelProtoToChannel(c *lnrpc.PendingChannelsResponse_ClosedChannel) *models.Channel {
	return &models.Channel{
		Status:        models.ChannelClosing,
		RemotePubKey:  c.Channel.RemoteNodePub,
		Capacity:      c.Channel.Capacity,
		LocalBalance:  c.Channel.LocalBalance,
		RemoteBalance: c.Channel.RemoteBalance,
		ChannelPoint:  c.Channel.ChannelPoint,
	}
}

func forceClosingChannelProtoToChannel(c *lnrpc.PendingChannelsResponse_ForceClosedChannel) *models.Channel {
	return &models.Channel{
		Status:            models.ChannelForceClosing,
		RemotePubKey:      c.Channel.RemoteNodePub,
		Capacity:          c.Channel.Capacity,
		LocalBalance:      c.Channel.LocalBalance,
		RemoteBalance:     c.Channel.RemoteBalance,
		ChannelPoint:      c.Channel.ChannelPoint,
		BlocksTilMaturity: c.BlocksTilMaturity,
	}
}

func waitingCloseChannelProtoToChannel(c *lnrpc.PendingChannelsResponse_WaitingCloseChannel) *models.Channel {
	return &models.Channel{
		Status:        models.ChannelWaitingClose,
		RemotePubKey:  c.Channel.RemoteNodePub,
		Capacity:      c.Channel.Capacity,
		LocalBalance:  c.Channel.LocalBalance,
		RemoteBalance: c.Channel.RemoteBalance,
		ChannelPoint:  c.Channel.ChannelPoint,
	}
}

func payreqProtoToPayReq(h *lnrpc.PayReq, payreq string) *models.PayReq {
	if h == nil {
		return nil
	}
	return &models.PayReq{
		Destination:     h.Destination,
		PaymentHash:     h.PaymentHash,
		Amount:          h.NumSatoshis,
		Timestamp:       h.Timestamp,
		Expiry:          h.Expiry,
		Description:     h.Description,
		DescriptionHash: h.DescriptionHash,
		FallbackAddr:    h.FallbackAddr,
		CltvExpiry:      h.CltvExpiry,
		String:          payreq,
	}
}

func sendPaymentProtoToPayment(payreq *models.PayReq, resp *lnrpc.SendResponse) *models.Payment {
	if payreq == nil || resp == nil {
		return nil
	}

	payment := &models.Payment{
		PaymentError:    resp.PaymentError,
		PaymentPreimage: resp.PaymentPreimage,
		PayReq:          payreq,
	}

	if resp.PaymentRoute != nil {
		payment.Route = &models.Route{
			TimeLock: resp.PaymentRoute.GetTotalTimeLock(),
			Fee:      resp.PaymentRoute.GetTotalFees(),
			Amount:   resp.PaymentRoute.GetTotalAmt(),
		}
	}

	return payment
}

func infoProtoToInfo(resp *lnrpc.GetInfoResponse) *models.Info {
	if resp == nil {
		return nil
	}

	chains := []string{}
	for i := range resp.Chains {
		chains = append(chains, resp.Chains[i].Chain)
	}

	return &models.Info{
		PubKey:              resp.IdentityPubkey,
		Alias:               resp.Alias,
		NumPendingChannels:  resp.NumPendingChannels,
		NumActiveChannels:   resp.NumActiveChannels,
		NumInactiveChannels: resp.NumInactiveChannels,
		NumPeers:            resp.NumPeers,
		BlockHeight:         resp.BlockHeight,
		BlockHash:           resp.BlockHash,
		Synced:              resp.SyncedToChain,
		Version:             resp.Version,
		Chains:              chains,
		Testnet:             resp.Testnet,
	}
}

func nodeProtoToNode(resp *lnrpc.NodeInfo) *models.Node {
	if resp == nil || resp.Node == nil {
		return nil
	}

	addresses := make([]*models.NodeAddress, len(resp.Node.Addresses))
	for i := range resp.Node.Addresses {
		addresses[i] = &models.NodeAddress{
			Network: resp.Node.Addresses[i].Network,
			Addr:    resp.Node.Addresses[i].Addr,
		}
	}
	channels := []*models.Channel{}
	for _, c := range resp.Channels {
		ch := &models.Channel{
			ID:           c.ChannelId,
			ChannelPoint: c.ChanPoint,
			Capacity:     c.Capacity,
			LocalPolicy:  protoToRoutingPolicy(c.Node1Policy),
			RemotePolicy: protoToRoutingPolicy(c.Node2Policy),
		}
		if c.Node1Pub != resp.Node.PubKey {
			ch.LocalPolicy, ch.RemotePolicy = ch.RemotePolicy, ch.LocalPolicy
		}
		channels = append(channels, ch)
	}

	return &models.Node{
		NumChannels:   resp.NumChannels,
		TotalCapacity: resp.TotalCapacity,
		LastUpdate:    time.Unix(int64(resp.Node.LastUpdate), 0),
		PubKey:        resp.Node.PubKey,
		Alias:         resp.Node.Alias,
		Addresses:     addresses,
		Channels:      channels,
	}
}

func protoToRoutingPolicy(resp *lnrpc.RoutingPolicy) *models.RoutingPolicy {
	if resp == nil {
		return nil
	}
	return &models.RoutingPolicy{
		TimeLockDelta:           resp.TimeLockDelta,
		MinHtlc:                 resp.MinHtlc,
		MaxHtlc:                 resp.MaxHtlcMsat,
		FeeBaseMsat:             resp.FeeBaseMsat,
		FeeRateMilliMsat:        resp.FeeRateMilliMsat,
		Disabled:                resp.Disabled,
		InboundFeeBaseMsat:      resp.InboundFeeBaseMsat,
		InboundFeeRateMilliMsat: resp.InboundFeeRateMilliMsat,
	}
}

func protoToTransactions(resp *lnrpc.TransactionDetails) []*models.Transaction {
	if resp == nil {
		return nil
	}

	transactions := make([]*models.Transaction, len(resp.Transactions))
	for i := range resp.Transactions {
		transactions[i] = protoToTransaction(resp.Transactions[i])
	}
	return transactions
}

func protoToTransaction(resp *lnrpc.Transaction) *models.Transaction {
	return &models.Transaction{
		TxHash:           resp.TxHash,
		Amount:           resp.Amount,
		NumConfirmations: resp.NumConfirmations,
		BlockHash:        resp.BlockHash,
		BlockHeight:      resp.BlockHeight,
		Date:             time.Unix(int64(resp.TimeStamp), 0),
		TotalFees:        resp.TotalFees,
		DestAddresses:    resp.DestAddresses,
	}
}

func protoToRoutingEvent(resp *routerrpc.HtlcEvent) *models.RoutingEvent {
	var status, direction int
	var incomingMsat, outgoingMsat uint64
	var incomingTimelock, outgoingTimelock uint32
	var amountMsat, feeMsat uint64
	var failureCode int32
	var detail string

	if fe := resp.GetForwardEvent(); fe != nil {
		status = models.RoutingStatusActive
		incomingMsat = fe.Info.IncomingAmtMsat
		outgoingMsat = fe.Info.OutgoingAmtMsat
		incomingTimelock = fe.Info.IncomingTimelock
		outgoingTimelock = fe.Info.OutgoingTimelock
	} else if ffe := resp.GetForwardFailEvent(); ffe != nil {
		status = models.RoutingStatusFailed
	} else if se := resp.GetSettleEvent(); se != nil {
		status = models.RoutingStatusSettled
	} else if lfe := resp.GetLinkFailEvent(); lfe != nil {
		incomingMsat = lfe.Info.IncomingAmtMsat
		outgoingMsat = lfe.Info.OutgoingAmtMsat
		incomingTimelock = lfe.Info.IncomingTimelock
		outgoingTimelock = lfe.Info.OutgoingTimelock
		status = models.RoutingStatusLinkFailed
		detail = lfe.WireFailure.String()
		if s := lfe.FailureDetail.String(); s != "" {
			detail = fmt.Sprintf("%s %s", detail, s)
		}
		if lfe.FailureString != "" {
			firstLine := strings.Split(lfe.FailureString, "\n")[0]
			detail = fmt.Sprintf("%s %s", detail, firstLine)
		}
		failureCode = int32(lfe.WireFailure)
	}

	switch resp.EventType {
	case routerrpc.HtlcEvent_SEND:
		direction = models.RoutingSend
		amountMsat = outgoingMsat
	case routerrpc.HtlcEvent_RECEIVE:
		direction = models.RoutingReceive
		amountMsat = incomingMsat
	case routerrpc.HtlcEvent_FORWARD:
		direction = models.RoutingForward
		amountMsat = outgoingMsat
		feeMsat = incomingMsat - outgoingMsat
	}

	return &models.RoutingEvent{
		IncomingChannelId: resp.IncomingChannelId,
		OutgoingChannelId: resp.OutgoingChannelId,
		IncomingHtlcId:    resp.IncomingHtlcId,
		OutgoingHtlcId:    resp.OutgoingHtlcId,
		LastUpdate:        time.Unix(0, int64(resp.TimestampNs)),
		Direction:         direction,
		Status:            status,
		IncomingTimelock:  incomingTimelock,
		OutgoingTimelock:  outgoingTimelock,
		AmountMsat:        amountMsat,
		FeeMsat:           feeMsat,
		FailureCode:       failureCode,
		FailureDetail:     detail,
	}
}

func protoToForwardingHistory(resp *lnrpc.ForwardingHistoryResponse) []*models.ForwardingEvent {
	if resp == nil {
		return nil
	}

	forwardingEvents := make([]*models.ForwardingEvent, len(resp.ForwardingEvents))
	for i := range resp.ForwardingEvents {
		forwardingEvents[i] = protoToForwardingEvent(resp.ForwardingEvents[i])
	}
	return forwardingEvents
}

func protoToForwardingEvent(resp *lnrpc.ForwardingEvent) *models.ForwardingEvent {
	return &models.ForwardingEvent{

		ChanIdIn:   resp.ChanIdIn,
		ChanIdOut:  resp.ChanIdOut,
		AmtIn:      resp.AmtIn,
		AmtOut:     resp.AmtOut,
		Fee:        resp.Fee,
		FeeMsat:    resp.FeeMsat,
		AmtInMsat:  resp.AmtInMsat,
		AmtOutMsat: resp.AmtOutMsat,
		EventTime:  time.Unix(0, int64(resp.TimestampNs)),
	}
}
