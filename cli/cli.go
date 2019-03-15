package cli

import (
	cli "gopkg.in/urfave/cli.v2"
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
		},
		Commands: []*cli.Command{},
	}
}
