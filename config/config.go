package config

import (
	"os"

	"github.com/dibyendu/Authentication-Authorization/lib/constants"
	"github.com/dibyendu/Authentication-Authorization/pkg/client/db"
	"github.com/dibyendu/Authentication-Authorization/pkg/client/redis"
)

type AppConfig struct {
	DB    *db.Config
	Redis *redis.Config
}

var (
	appConfig AppConfig
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
		Redis: &redis.Config{
			Host:     os.Getenv("REDIS_HOST"),
			Port:     os.Getenv("REDIS_PORT"),
			Password: os.Getenv("REDIS_PASS"),
		},
	}
	return &appConfig
}
