module github.com/edouardparis/lntop

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/fatih/color v1.7.0
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/jroimartin/gocui v0.4.0
	github.com/juju/clock v0.0.0-20190205081909-9c5c9712527c // indirect
	github.com/lightningnetwork/lnd v0.13.2-beta
	github.com/mattn/go-colorable v0.1.1 // indirect
	github.com/mattn/go-isatty v0.0.7 // indirect
	github.com/miekg/dns v1.1.6 // indirect
	github.com/nsf/termbox-go v0.0.0-20190121233118-02980233997d // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/tsdb v0.7.1 // indirect
	go.uber.org/zap v1.14.1
	golang.org/x/text v0.3.3
	google.golang.org/grpc v1.29.1
	gopkg.in/macaroon-bakery.v2 v2.1.0 // indirect
	gopkg.in/macaroon.v2 v2.1.0
	gopkg.in/resty.v1 v1.12.0 // indirect
	gopkg.in/urfave/cli.v2 v2.0.0-20180128182452-d3ae77c26ac8
)

replace go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20201125193152-8a03d2e9614b
