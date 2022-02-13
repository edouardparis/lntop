package config

import (
	"fmt"
	"os"
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
columns = [
	"STATUS",      # status of the channel
	"ALIAS",       # alias of the channel node
	"GAUGE",       # ascii bar with percent local/capacity
	"LOCAL",       # the local amount of the channel
	# "REMOTE",    # the remote amount of the channel
	"CAP",         # the total capacity of the channel
	"SENT",        # the total amount sent
	"RECEIVED",    # the total amount received
	"HTLC",        # the number of pending HTLC
	"UNSETTLED",   # the amount unsettled in the channel
	"CFEE",        # the commit fee
	"LAST UPDATE", # last update of the channel
	"PRIVATE",     # true if channel is private
	"ID",          # the id of the channel
	# "SCID",      # short channel id (BxTxO formatted)
	# "NUPD",      # number of channel updates
]

[views.transactions]
# It is possible to add, remove and order columns of the
# table with the array columns. The available values are:
columns = [
	"DATE",      # date of the transaction
	"HEIGHT",    # block height of the transaction
	"CONFIR",    # number of confirmations
	"AMOUNT",    # amount moved by the transaction
	"FEE",       # fee of the transaction
	"ADDRESSES", # number of transaction output addresses
]

[views.routing]
columns = [
	"DIR",            # event type:  send, receive, forward
	"STATUS",         # one of: active, settled, failed, linkfail
	"IN_CHANNEL",     # channel id of the incomming channel
	"IN_ALIAS",       # incoming channel node alias
	# "IN_SCID",      # incoming short channel id (BxTxO)
	# "IN_HTLC",      # htlc id on incoming channel
	# "IN_TIMELOCK",  # incoming timelock height
	"OUT_CHANNEL",    # channel id of the outgoing channel
	"OUT_ALIAS",      # outgoing channel node alias
	# "OUT_SCID",     # outgoing short channel id (BxTxO)
	# "OUT_HTLC",     # htlc id on outgoing channel
	# "OUT_TIMELOCK", # outgoing timelock height
	"AMOUNT",         # routed amount
	"FEE",            # routing fee
	"LAST UPDATE",    # last update
	"DETAIL",         # error description
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
	lndAddress, present := os.LookupEnv("LND_ADDRESS")
	if !present {
		lndAddress = "//127.0.0.1:10009"
	}
	certPath, present := os.LookupEnv("CERT_PATH")
	if !present {
		certPath = path.Join(usr.HomeDir, ".lnd/tls.cert")
	}
	macaroonPath, present := os.LookupEnv("MACAROON_PATH")
	if !present {
		macaroonPath = path.Join(usr.HomeDir, ".lnd/data/chain/bitcoin/mainnet/readonly.macaroon")
	}
	return &Config{
		Logger: Logger{
			Type: "production",
			Dest: path.Join(usr.HomeDir, ".lntop/lntop.log"),
		},
		Network: Network{
			Name:            "lnd",
			Type:            "lnd",
			Address:         lndAddress,
			Cert:            certPath,
			Macaroon:        macaroonPath,
			MacaroonTimeOut: 60,
			MaxMsgRecvSize:  52428800,
			ConnTimeout:     1000000,
			PoolCapacity:    4,
		},
	}
}
