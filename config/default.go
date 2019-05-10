package config

import (
	"fmt"
	"os/user"
	"path"
)

func DefaultFileContent() string {
	cfg := NewDefault()
	return fmt.Sprintf(`
[logger]
type = "%[1]s"
dest = "%[2]s"

[network]
name = "%[3]s"
type = "%[4]s"
address = "%[5]s"
cert = "%[6]s"
macaroon = "%[7]s"
macaroon_timeout = %[8]d
max_msg_recv_size = %[9]d
conn_timeout = %[10]d
pool_capacity = %[11]d

[views]
# views.channels is the view displaying channel list.
[views.channels]
# It is possible to add, remove and order columns of the
# table with the array columns. The available values are:
# STATUS      status of the channel
# ALIAS       alias of the channel node
# GAUGE       ascii bar with percent local/capacity
# LOCAL       the local amount of the channel
# CAP         the total capacity of the channel
# HTLC        the number of pending HTLC
# UNSETTLED   the amount unsettled in the channel
# CFEE        the commit fee
# LAST UPDATE last update of the channel
# PRIVATE     true if channel is private
# ID          the id of the channel

columns = [
	"STATUS",
	"ALIAS",
	"GAUGE",
	"LOCAL",
	"CAP",
	"HTLC",
	"UNSETTLED",
	"CFEE",
	"LAST UPDATE",
	"PRIVATE",
	"ID",
]

[views.transactions]
# It is possible to add, remove and order columns of the
# table with the array columns. The available values are:
# DATE      date of the transaction
# HEIGHT    block height of the transaction
# CONFIR    number of confirmations
# AMOUNT    amount moved by the transaction
# FEE       fee of the transaction
# ADDRESSES number of transaction output addresses

columns = [
	"TIME",
	"HEIGHT",
	"CONFIR",
	"AMOUNT",
	"FEE",
	"ADDRESSES",
]
`,
		cfg.Logger.Type,
		cfg.Logger.Dest,
		cfg.Network.Name,
		cfg.Network.Type,
		cfg.Network.Address,
		cfg.Network.Cert,
		cfg.Network.Macaroon,
		cfg.Network.MacaroonTimeOut,
		cfg.Network.MaxMsgRecvSize,
		cfg.Network.ConnTimeout,
		cfg.Network.PoolCapacity,
	)
}

func NewDefault() *Config {
	usr, _ := user.Current()
	return &Config{
		Logger: Logger{
			Type: "production",
			Dest: path.Join(usr.HomeDir, ".lntop/lntop.log"),
		},
		Network: Network{
			Name:            "lnd",
			Type:            "lnd",
			Address:         "//127.0.0.1:10009",
			Cert:            path.Join(usr.HomeDir, ".lnd/tls.cert"),
			Macaroon:        path.Join(usr.HomeDir, ".lnd/data/chain/bitcoin/mainnet/admin.macaroon"),
			MacaroonTimeOut: 60,
			MaxMsgRecvSize:  52428800,
			ConnTimeout:     1000000,
			PoolCapacity:    3,
		},
	}
}
