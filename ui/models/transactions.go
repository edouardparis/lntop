package models

import (
	"context"

	"github.com/edouardparis/lntop/network/models"
)

type Transactions struct {
	list []*models.Transaction
}

func (t Transactions) List() []*models.Transaction {
	return t.list
}

func (m *Models) RefreshTransactions(ctx context.Context) error {
	transactions, err := m.network.GetTransactions(ctx)
	if err != nil {
		return err
	}
	*m.Transactions = Transactions{list: transactions}
	return nil
}
