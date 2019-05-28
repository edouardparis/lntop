package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	cli "gopkg.in/urfave/cli.v2"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/pubsub"
	"github.com/edouardparis/lntop/ui"
	"github.com/edouardparis/lntop/version"
)

// newApp creates a new cli app.
func newApp() *cli.App {
	cli.VersionFlag = &cli.BoolFlag{
		Name: "version", Aliases: []string{},
		Usage: "print the version",
	}

	return &cli.App{
		Name:                  "lntop",
		Version:               version.Version,
		EnableShellCompletion: true,
		Action:                run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "path to config file",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "pubsub",
				Aliases: []string{""},
				Usage:   "run the pubsub only",
				Action:  pubsubRun,
			},
		},
	}
}

func run(c *cli.Context) error {
	cfg, err := config.Load(c.String("config"))
	if err != nil {
		return err
	}

	app, err := app.New(cfg)
	if err != nil {
		return err
	}

	ctx := context.Background()

	ps := pubsub.New(app.Logger, app.Network)
	go ps.Run(ctx)

	err = ui.Run(ctx, app, ps)
	if err != nil {
		app.Logger.Debug("ui", logging.String("error", err.Error()))
	}

	return ps.Stop()
}

func pubsubRun(c *cli.Context) error {
	cfg, err := config.Load(c.String("config"))
	if err != nil {
		return err
	}

	app, err := app.New(cfg)
	if err != nil {
		return err
	}

	ps := pubsub.New(app.Logger, app.Network)
	ps.Run(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	go func() {
		<-sig
		ps.Stop()
	}()

	return nil
}

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
