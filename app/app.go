package app

import (
	"github.com/edouardparis/lntop/config"
	"github.com/edouardparis/lntop/logging"
	"github.com/edouardparis/lntop/network"
)

type App struct {
	Config  *config.Config
	Logger  logging.Logger
	Network *network.Network
}

func New(cfg *config.Config) (*App, error) {
	logger := logging.New(cfg.Logger)

	network, err := network.New(&cfg.Network, logger)
	if err != nil {
		return nil, err
	}

	return &App{
		Config:  cfg,
		Logger:  logger,
		Network: network,
	}, nil
}
