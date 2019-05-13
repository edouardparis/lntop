package models

import (
	"sort"
	"sync"

	"github.com/edouardparis/lntop/network/models"
)

type ChannelsSort func(*models.Channel, *models.Channel) bool

type Channels struct {
	current *models.Channel
	index   map[string]*models.Channel
	list    []*models.Channel
	sort    ChannelsSort
	mu      sync.RWMutex
}

func (c *Channels) List() []*models.Channel {
	return c.list
}

func (c *Channels) Len() int {
	return len(c.list)
}

func (c *Channels) Swap(i, j int) {
	c.list[i], c.list[j] = c.list[j], c.list[i]
}

func (c *Channels) Less(i, j int) bool {
	return c.sort(c.list[i], c.list[j])
}

func (c *Channels) Sort(s ChannelsSort) {
	if s == nil {
		return
	}
	c.sort = s
	sort.Sort(c)
}

func (c *Channels) Current() *models.Channel {
	return c.current
}

func (c *Channels) SetCurrent(index int) {
	c.current = c.Get(index)
}

func (c *Channels) Get(index int) *models.Channel {
	if index < 0 || index > len(c.list)-1 {
		return nil
	}

	return c.list[index]
}

func (c *Channels) GetByChanPoint(chanPoint string) *models.Channel {
	return c.index[chanPoint]
}

func (c *Channels) Contains(channel *models.Channel) bool {
	_, ok := c.index[channel.ChannelPoint]
	return ok
}

func (c *Channels) Add(channel *models.Channel) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Contains(channel) {
		return
	}
	c.index[channel.ChannelPoint] = channel
	c.list = append(c.list, channel)
}

func (c *Channels) Update(newChannel *models.Channel) {
	c.mu.Lock()
	defer c.mu.Unlock()

	oldChannel, ok := c.index[newChannel.ChannelPoint]
	if !ok {
		c.Add(newChannel)
		if c.sort != nil {
			sort.Sort(c)
		}
		return
	}

	oldChannel.ID = newChannel.ID
	oldChannel.Status = newChannel.Status
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

	if newChannel.LastUpdate != nil {
		oldChannel.LastUpdate = newChannel.LastUpdate
	}

	if newChannel.Policy1 != nil {
		oldChannel.Policy1 = newChannel.Policy1
	}

	if newChannel.Policy2 != nil {
		oldChannel.Policy2 = newChannel.Policy2
	}
}

func NewChannels() *Channels {
	return &Channels{
		list:  []*models.Channel{},
		index: make(map[string]*models.Channel),
	}
}
