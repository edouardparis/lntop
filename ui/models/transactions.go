package models

import (
	"context"

	"github.com/edouardparis/lntop/network/models"
)

type Transactions struct {
	current int
	list    []*models.Transaction
}

func (t Transactions) Current() *models.Transaction {
	return t.Get(t.current)
}

func (t Transactions) List() []*models.Transaction {
	return t.list
}

func (t *Transactions) Len() int {
	return len(t.list)
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
