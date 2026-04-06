package app

import (
	"github.com/adishjain1107/tradex/pkg/auth/config"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Config *config.Config
	DB     *mongo.Database
}

func New(cfg *config.Config, db *mongo.Database) *App {
	return &App{
		Config: cfg,
		DB:     db,
	}
}
