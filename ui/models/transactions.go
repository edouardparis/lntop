package models

import (
	"context"
	"sort"
	"sync"

	"github.com/edouardparis/lntop/network/models"
)

type TransactionsSort func(*models.Transaction, *models.Transaction) bool

type Transactions struct {
	current *models.Transaction
	list    []*models.Transaction
	sort    TransactionsSort
	mu      sync.RWMutex
}

func (t *Transactions) Current() *models.Transaction {
	return t.current
}

func (t *Transactions) SetCurrent(index int) {
	t.current = t.Get(index)
}

func (t *Transactions) List() []*models.Transaction {
	return t.list
}

func (t *Transactions) Len() int {
	return len(t.list)
}

func (t *Transactions) Swap(i, j int) {
	t.list[i], t.list[j] = t.list[j], t.list[i]
}

func (t *Transactions) Less(i, j int) bool {
	return t.sort(t.list[i], t.list[j])
}

func (t *Transactions) Sort(s TransactionsSort) {
	if s == nil {
		return
	}
	t.sort = s
	sort.Sort(t)
}

func (t *Transactions) Get(index int) *models.Transaction {
	if index < 0 || index > len(t.list)-1 {
		return nil
	}

	return t.list[index]
}

func (t *Transactions) Contains(tx *models.Transaction) bool {
	if tx == nil {
		return false
	}
	for i := range t.list {
		if t.list[i].TxHash == tx.TxHash {
			return true
		}
	}
	return false
}

func (t *Transactions) Add(tx *models.Transaction) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.Contains(tx) {
		return
	}
	t.list = append(t.list, tx)
	if t.sort != nil {
		sort.Sort(t)
	}
}

func (t *Transactions) Update(tx *models.Transaction) {
	if tx == nil {
		return
	}
	if !t.Contains(tx) {
		t.Add(tx)
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	for i := range t.list {
		if t.list[i].TxHash == tx.TxHash {
			t.list[i].NumConfirmations = tx.NumConfirmations
			t.list[i].BlockHeight = tx.BlockHeight
		}
	}

	if t.sort != nil {
		sort.Sort(t)
	}
}

func (m *Models) RefreshTransactions(ctx context.Context) error {
	transactions, err := m.network.GetTransactions(ctx)
	if err != nil {
		return err
	}

	for i := range transactions {
		m.Transactions.Update(transactions[i])
	}

	return nil
}
