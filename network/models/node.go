package models

import "time"

type Node struct {
	NumChannels   uint32
	TotalCapacity int64
	LastUpdate    time.Time
	PubKey        string
	Alias         string
	ForcedAlias   string
	Addresses     []*NodeAddress
}

type NodeAddress struct {
	Network string
	Addr    string
}
