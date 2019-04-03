package models

import "github.com/edouardparis/lntop/network/models"

type Channels struct {
	index map[uint64]*models.Channel
	list  []*models.Channel
}

func (c Channels) List() []*models.Channel {
	return c.list
}

func (c *Channels) Get(index int) *models.Channel {
	if index < 0 || index > len(c.list)-1 {
		return nil
	}

	return c.list[index]
}

func (c *Channels) GetByID(id uint64) *models.Channel {
	return c.index[id]
}

func (c Channels) Contains(channel *models.Channel) bool {
	_, ok := c.index[channel.ID]
	return ok
}

func (c *Channels) Add(channel *models.Channel) {
	if c.Contains(channel) {
		return
	}
	c.index[channel.ID] = channel
	c.list = append(c.list, channel)
}

func (c *Channels) Update(newChannel *models.Channel) {
	oldChannel, ok := c.index[newChannel.ID]
	if !ok {
		c.Add(newChannel)
		return
	}
	oldChannel.Active = newChannel.Active
	oldChannel.LocalBalance = newChannel.LocalBalance
	oldChannel.RemoteBalance = newChannel.RemoteBalance
	oldChannel.CommitFee = newChannel.CommitFee
	oldChannel.CommitWeight = newChannel.CommitWeight
	oldChannel.FeePerKiloWeight = newChannel.FeePerKiloWeight
	oldChannel.UnsettledBalance = newChannel.UnsettledBalance
	oldChannel.TotalAmountSent = newChannel.TotalAmountSent
	oldChannel.TotalAmountReceived = newChannel.TotalAmountReceived
	oldChannel.UpdatesCount = newChannel.UpdatesCount
	oldChannel.CSVDelay = newChannel.CSVDelay
	oldChannel.Private = newChannel.Private
	oldChannel.PendingHTLC = newChannel.PendingHTLC

	if newChannel.LastUpdated != nil {
		oldChannel.LastUpdated = newChannel.LastUpdated
	}
}

func NewChannels() *Channels {
	return &Channels{
		list:  []*models.Channel{},
		index: make(map[uint64]*models.Channel),
	}
}

type Channel struct {
	Item *models.Channel
}
