package cli

import (
	"context"
	"fmt"

	cli "gopkg.in/urfave/cli.v2"

	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network"
	"github.com/edouardparis/lntop/network/backend"
)

// New creates a new cli app.
func New() *cli.App {
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{},
		Usage: "print the version",
	}

	return &cli.App{
		Name:                  "lntop",
		EnableShellCompletion: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "v",
				Usage: "verbose",
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "verbose",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "wallet-balance",
				Aliases: []string{""},
				Usage:   "",
				Action:  walletBalance,
			},
		},
	}
}

func getNetworkFromConfig(c *cli.Context) (backend.Backend, error) {
	cfg, err := config.Load(c.String("config"))
	if err != nil {
		return nil, err
	}

	logger := logging.New(config.Logger{Type: "development"})

	return network.New(&cfg.Network, logger)
}

func walletBalance(c *cli.Context) error {
	clt, err := getNetworkFromConfig(c)
	if err != nil {
		return err
	}

	res, err := clt.GetWalletBalance(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(res.TotalBalance)

	return nil
}
