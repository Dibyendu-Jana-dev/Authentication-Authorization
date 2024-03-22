package main

import (
	"github.com/dibyendu/Authentication-Authorization/config"
	_ "github.com/dibyendu/Authentication-Authorization/docs"
	"github.com/dibyendu/Authentication-Authorization/pkg/app"
	"github.com/joho/godotenv"
	"log"
)

// @title Authentication and Authorization API
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @securityDefinitions.apikey bearerToken
// @in header
// @name Authorization

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @description This is the API documentation for the Authentication and Authorization service.
// @version 2.0
// @host localhost:8080
// @BasePath /

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file: " + err.Error())
	}
	config := config.Init()
	app.StartApp(config) 
}
