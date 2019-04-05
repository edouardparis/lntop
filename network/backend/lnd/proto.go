package lnd

import (
	"time"

	"github.com/lightningnetwork/lnd/lnrpc"

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
	return &models.Invoice{
		Expiry:         req.GetExpiry(),
		Amount:         req.GetValue(),
		Description:    req.GetMemo(),
		CreationDate:   req.GetCreationDate(),
		RHash:          resp.GetRHash(),
		PaymentRequest: resp.GetPaymentRequest(),
		Index:          resp.GetAddIndex(),
	}
}

func lookupInvoiceProtoToInvoice(resp *lnrpc.Invoice) *models.Invoice {
	return &models.Invoice{
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
	return &models.Channel{
		ID:                  c.GetChanId(),
		Active:              c.GetActive(),
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
		Chains:              resp.Chains,
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

	return &models.Node{
		NumChannels:   resp.NumChannels,
		TotalCapacity: resp.TotalCapacity,
		LastUpdate:    time.Unix(int64(resp.Node.LastUpdate), 0),
		PubKey:        resp.Node.PubKey,
		Alias:         resp.Node.Alias,
		Addresses:     addresses,
	}
}

func protoToRoutingPolicy(resp *lnrpc.RoutingPolicy) *models.RoutingPolicy {
	if resp == nil {
		return nil
	}
	return &models.RoutingPolicy{
		TimeLockDelta:    resp.TimeLockDelta,
		MinHtlc:          resp.MinHtlc,
		FeeBaseMsat:      resp.FeeBaseMsat,
		FeeRateMilliMsat: resp.FeeRateMilliMsat,
		Disabled:         resp.Disabled,
	}
}
