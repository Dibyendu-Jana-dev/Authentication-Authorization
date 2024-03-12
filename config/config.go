package config

import (
	"github.com/dibyendu/Authentication-Authorization/lib/constants"
	"github.com/dibyendu/Authentication-Authorization/pkg/client/db"
	"os"
)

type AppConfig struct {
	DB       *db.Config
}

var (
	appConfig     AppConfig
)

func Init() *AppConfig {
	userCollection := make(map[string]string)
	userCollection["user"] = constants.USER_COLLECTION

	appConfig = AppConfig{
		DB: &db.Config{
			Host:           os.Getenv("DB_HOST"),
			Port:           os.Getenv("DB_PORT"),
			MaxPool:        os.Getenv("MAX_POOL"),
			Database:       os.Getenv("DB_NAME"),
			UserCollection: userCollection,
		},
	}
	return &appConfig
}
