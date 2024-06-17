package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)
const DriverName = "mongodb"
type Config struct {
	Host           string
	Port           string
	MaxPool        string
	Database       string
	UserCollection map[string]string
}

func Init(dbConfig *Config) (*mongo.Client, error) {
	dataSource := fmt.Sprintf("%s://%s:%s/?maxPoolSize=%s&w=majority",
		DriverName, dbConfig.Host, dbConfig.Port, dbConfig.MaxPool)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dataSource))
	if err != nil {
		return nil, err
	}
	// verifies connection is db is working
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}
	return client, nil
}
