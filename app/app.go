package app

import (
	"github.com/adishjain1107/tradex/pkg/auth/config"
)

type Application struct {
	Config *config.Config
}


func NewApp(cfg *config.Config) *Application {
	return &Application{
		Config: cfg,
	}
}