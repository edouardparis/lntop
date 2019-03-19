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

func Load() (*App, error) {
	return &App{}, nil
}
