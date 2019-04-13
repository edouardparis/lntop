package cli

import (
	"context"

	cli "gopkg.in/urfave/cli.v2"

	"github.com/edouardparis/lntop/app"
	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/events"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/pubsub"
	"github.com/edouardparis/lntop/ui"
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

	events := make(chan *events.Event)

	go func() {
		err := ui.Run(context.Background(), app, events)
		if err != nil {
			app.Logger.Debug("ui", logging.String("error", err.Error()))
		}
	}()

	pubsub.Run(context.Background(), app, events)
	return nil
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

	events := make(chan *events.Event)
	pubsub.Run(context.Background(), app, events)
	//ev := <-events
	//app.Logger.Info("events quit ", logging.String("type", ev.Type))

	return nil
}
