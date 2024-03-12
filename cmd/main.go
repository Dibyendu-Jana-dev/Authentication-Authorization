package main

import (
	"github.com/dibyendu/Authentication-Authorization/config"
	"github.com/dibyendu/Authentication-Authorization/pkg/app"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load("dev.env")
	if err != nil {
		log.Println("Error loading .env file: " + err.Error())
	}
	config := config.Init()
	app.StartApp(config) 
}
