package app

import "github.com/adishjain1107/tradex/pkg/order-matcher/config"

type App struct {
	Config *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		Config: cfg,
	}
}
