package config

type Network struct {
	ID              string
	Type            string
	Status          string
	Address         string
	Cert            string
	Macaroon        string
	MacaroonTimeOut int64
	MacaroonIP      string
	MaxMsgRecvSize  int
	ConnTimeout     int
	PoolCapacity    int
}
