package models

import (
	"context"

	"github.com/edouardparis/lntop/network/models"
)

type TransactionsSort func(*models.Transaction, *models.Transaction) bool

type Transactions struct {
	current *models.Transaction
	list    []*models.Transaction
	sort    TransactionsSort
}

func (t Transactions) Current() *models.Transaction {
	return t.current
}

func (t *Transactions) SetCurrent(index int) {
	t.current = t.Get(index)
}

func (t Transactions) List() []*models.Transaction {
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

func (t *Transactions) WithSort(s TransactionsSort) {
	t.sort = s
}

func (t *Transactions) Get(index int) *models.Transaction {
	if index < 0 || index > len(t.list)-1 {
		return nil
	}

	return t.list[index]
}

func (m *Models) RefreshTransactions(ctx context.Context) error {
	transactions, err := m.network.GetTransactions(ctx)
	if err != nil {
		return err
	}
	*m.Transactions = Transactions{list: transactions}
	return nil
}
