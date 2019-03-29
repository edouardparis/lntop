package models

import "github.com/edouardparis/lntop/logging"

type Info struct {
	PubKey              string
	Alias               string
	NumPendingChannels  uint32
	NumActiveChannels   uint32
	NumInactiveChannels uint32
	NumPeers            uint32
	BlockHeight         uint32
	BlockHash           string
	Synced              bool
	Version             string
	Chains              []string
	Testnet             bool
}

func (i Info) MarshalLogObject(enc logging.ObjectEncoder) error {
	enc.AddString("pubkey", i.PubKey)
	enc.AddString("Alias", i.Alias)
	return nil
}
