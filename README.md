# lntop

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/edouardparis/lntop/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/edouardparis/lntop)](https://goreportcard.com/report/github.com/edouardparis/lntop)
[![Godoc](https://godoc.org/github.com/edouardparis/lntop?status.svg)](https://godoc.org/github.com/edouardparis/lntop)
[![tippin.me](https://badgen.net/badge/%E2%9A%A1%EF%B8%8Ftippin.me/@edouardparis/F0918E)](https://tippin.me/@edouardparis)

`lntop` is an interactive text-mode channels viewer for Unix systems.

 ![lntop-v0.1.0](http://paris.iiens.net/lntop-v0.1.0.png)
 *lntop-v0.1.0*

## Install

Require the [go programming language](https://golang.org/) (version >= 1.11)
```
git clone https://github.com/edouardparis/lntop.git
cd lntop && export GO111MODULE=on && go install ./...
```

## Config

First time `lntop` is used a config file `.lntop/config.toml` is created
in the user home directory.

Change macaroon path according to your network.

```toml
[logger]
type = "production"
dest = "/root/.lntop/lntop.log"

[network]
name = "lnd"
type = "lnd"
address = "//127.0.0.1:10009"
cert = "/root/.lnd/tls.cert"
macaroon = "/root/.lnd/data/chain/bitcoin/mainnet/admin.macaroon"
macaroon_timeout = 60
max_msg_recv_size = 52428800
conn_timeout = 1000000
pool_capacity = 4

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
	"CAP",         # the total capacity of the channel
	"SENT",        # the total amount sent
	"RECEIVED",    # the total amount received
	"HTLC",        # the number of pending HTLC
	"UNSETTLED",   # the amount unsettled in the channel
	"CFEE",        # the commit fee
	"LAST UPDATE", # last update of the channel
	"PRIVATE",     # true if channel is private
	"ID",          # the id of the channel
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
	# "IN_HTLC",      # htlc id on incoming channel
	# "IN_TIMELOCK",  # incoming timelock height
	"OUT_CHANNEL",    # channel id of the outgoing channel
	"OUT_ALIAS",      # outgoing channel node alias
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

## Compatibility

| lntop  | lightningnetwork/lnd |
|--------|----------------------|
| v0.0.1 | v0.5.1               |
