package app

import "github.com/adishjain1107/tradex/pkg/notifier/config"

type App struct {
	Config *config.Config
}

func New(cfg *config.Config) *App {
	return &App{
		Config: cfg,
	}
}
