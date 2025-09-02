package models

import (
	"context"
	"sort"
	"sync"

	netmodels "github.com/edouardparis/lntop/network/models"
)

type ReceivedSort func(*netmodels.Invoice, *netmodels.Invoice) bool

type Received struct {
	list          []*netmodels.Invoice
	sort          ReceivedSort
	mu            sync.RWMutex
	StartDateUnix int64
}

func (r *Received) List() []*netmodels.Invoice { return r.list }
func (r *Received) Len() int                   { return len(r.list) }
func (r *Received) Swap(i, j int)              { r.list[i], r.list[j] = r.list[j], r.list[i] }
func (r *Received) Less(i, j int) bool         { return r.sort(r.list[i], r.list[j]) }

func (r *Received) Sort(s ReceivedSort) {
	if s == nil {
		return
	}
	r.sort = s
	sort.Sort(r)
}

func (r *Received) Contains(inv *netmodels.Invoice) bool {
	if inv == nil {
		return false
	}
	for _, it := range r.list {
		if it.GetRHash() == inv.GetRHash() {
			return true
		}
	}
	return false
}

func (r *Received) Add(inv *netmodels.Invoice) {
	if inv == nil {
		return
	}
	if !inv.Settled {
		return
	}
	// Apply start date filter if set
	if r.StartDateUnix > 0 {
		ts := inv.SettleDate
		if ts == 0 {
			ts = inv.CreationDate
		}
		if ts < r.StartDateUnix {
			return
		}
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Contains(inv) {
		return
	}
	r.list = append(r.list, inv)
	if r.sort != nil {
		sort.Sort(r)
	}
}

// RefreshReceived consumes an update (expected *netmodels.Invoice) and updates the list.
func (m *Models) RefreshReceived(update interface{}) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		inv, ok := update.(*netmodels.Invoice)
		if !ok {
			return nil
		}
		m.Received.Add(inv)
		return nil
	}
}

// RefreshReceivedFromNetwork fetches invoices from backend and populates Received model.
func (m *Models) RefreshReceivedFromNetwork(ctx context.Context) error {
	invoices, err := m.network.ListInvoices(ctx)
	if err != nil {
		return err
	}
	for _, inv := range invoices {
		if inv == nil || !inv.Settled {
			continue
		}
		// Add() will apply start date filter as needed
		m.Received.Add(inv)
	}
	return nil
}
