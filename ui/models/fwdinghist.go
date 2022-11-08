package models

import (
	"sort"
	"sync"

	"github.com/edouardparis/lntop/network/models"
)

type FwdinghistSort func(*models.ForwardingEvent, *models.ForwardingEvent) bool

type FwdingHist struct {
	StartTime    string
	MaxNumEvents uint32
	current      *models.ForwardingEvent
	list         []*models.ForwardingEvent
	sort         FwdinghistSort
	mu           sync.RWMutex
}

func (t *FwdingHist) Current() *models.ForwardingEvent {
	return t.current
}

func (t *FwdingHist) SetCurrent(index int) {
	t.current = t.Get(index)
}

func (t *FwdingHist) List() []*models.ForwardingEvent {
	return t.list
}

func (t *FwdingHist) Len() int {
	return len(t.list)
}

func (t *FwdingHist) Clear() {
	t.list = []*models.ForwardingEvent{}
}

func (t *FwdingHist) Swap(i, j int) {
	t.list[i], t.list[j] = t.list[j], t.list[i]
}

func (t *FwdingHist) Less(i, j int) bool {
	return t.sort(t.list[i], t.list[j])
}

func (t *FwdingHist) Sort(s FwdinghistSort) {
	if s == nil {
		return
	}
	t.sort = s
	sort.Sort(t)
}

func (t *FwdingHist) Get(index int) *models.ForwardingEvent {
	if index < 0 || index > len(t.list)-1 {
		return nil
	}

	return t.list[index]
}

func (t *FwdingHist) Update(events []*models.ForwardingEvent) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Clear()
	for _, event := range events {
		t.list = append(t.list, event)
	}
}
