# lntop

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/edouardparis/lntop/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/edouardparis/lntop)](https://goreportcard.com/report/github.com/edouardparis/lntop)
[![Godoc](https://godoc.org/github.com/edouardparis/lntop?status.svg)](https://godoc.org/github.com/edouardparis/lntop)

`lntop` is an interactive text-mode channels viewer for Unix systems.

 <img src="lntop-v0.1.0.png">
 *lntop-v0.1.0*

## Install

Require the [go programming language](https://golang.org/) (version >= 1.13)

**Raspberry Pi users: be aware that Raspbian ships with Go 1.11** ( see
[#30](https://github.com/edouardparis/lntop/issues/30) )

```
git clone https://github.com/edouardparis/lntop.git
cd lntop
go build   // creates a binary `lntop` in directory
go install // creates a binary and move it in your $GOBIN path

```

With Go version >= 1.16, you can use [`go-install`](https://golang.org/ref/mod#go-install)

```
go install github.com/edouardparis/lntop@latest
```

Note: If you are using [**Umbrel**](https://getumbrel.com) or [**Citadel**](https://runcitadel.space) you can simply install the [**Lightning Shell**](https://lightningshell.app) app from the respective dashboard. This will give you `lntop` which should just work without any additional configuration.

## Config

First time `lntop` is started, a config file `.lntop/config.toml` is created in the user's home directory. Change `address`, `cert` path and `macaroon` path according to your setup.

The following environment variables, if present, will be used in the initial config file instead of the defaults, so you won't have to have `lntop` fail on the first start and then manually edit the config file: `LND_ADDRESS`, `CERT_PATH`, `MACAROON_PATH`.

```toml
[logger]
type = "production"
dest = "/root/.lntop/lntop.log"

[network]
name = "lnd"
type = "lnd"
address = "//127.0.0.1:10009"
cert = "/root/.lnd/tls.cert"
macaroon = "/root/.lnd/data/chain/bitcoin/mainnet/readonly.macaroon"
macaroon_timeout = 60
max_msg_recv_size = 52428800
conn_timeout = 1000000
pool_capacity = 4

[network.aliases]
# Not all peers have aliases set up. In order to remember who is whom, pubkeys can be annotated.
# "Forced" aliases will be printed in a different color to be distinguished from network advertised aliases.
035e4ff418fc8b5554c5d9eea66396c227bd429a3251c8cbc711002ba215bfc226 = "Wallet of Satoshi"
03864ef025fde8fb587d989186ce6a4a186895ee44a926bfc370e2c366597a3f8f = "-=[ACINQ]=-"

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
	"REMOTE",    # the remote amount of the channel
	#"BASE_OUT"    # the outgoing base fee of the channel
	#"RATE_OUT"    # the outgoing fee rate in ppm of the channel
	#"BASE_IN"    # the incoming base fee of the channel
	#"RATE_IN"    # the incoming fee rate in ppm of the channel
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
```

## Routing view

Routing view displays screenful of latest routing events. This information
is not persisted in LND so the view always starts empty and is lost once
you exit `lntop`.

The events are in one of four states:

* `active` - HTLC pending
* `settled` - preimage revealed, HTLC removed
* `failed` - payment failed at a downstream node
* `linkfail` - payment failed at this node

## Docker

If you prefer to run `lntop` from a docker container, `cd docker` and follow [`README`](docker/README.md) there.
