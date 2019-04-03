package models

import "github.com/edouardparis/lntop/network/models"

type Channels struct {
	Items []*models.Channel
}

func (c *Channels) Get(index int) *models.Channel {
	if index < 0 || index > len(c.Items)-1 {
		return nil
	}

	return c.Items[index]
}

type Channel struct {
	Item *models.Channel
}
