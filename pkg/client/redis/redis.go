package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     string
	Password string
}

func Init(redConfig *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redConfig.Host + ":" + redConfig.Port,
		Password: redConfig.Password, //no password set
		DB:       0,  //use default DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil{
		log.Println("Error connecting")
		panic(err)
	}
	return client, nil
}
