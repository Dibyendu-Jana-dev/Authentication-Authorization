package app

import (
	"log"
	"net/http"

	"github.com/dibyendu/Authentication-Authorization/config"
	"github.com/dibyendu/Authentication-Authorization/pkg/client/db"
	"github.com/dibyendu/Authentication-Authorization/pkg/client/redis"
	"github.com/dibyendu/Authentication-Authorization/pkg/domain"
	"github.com/dibyendu/Authentication-Authorization/pkg/handler"
	"github.com/dibyendu/Authentication-Authorization/pkg/middleware"
	"github.com/dibyendu/Authentication-Authorization/pkg/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func StartApp(config *config.AppConfig) {
	log.Println("Starting app")

	// Initialize the database
	dbClient, err := db.Init(config.DB)
	if err != nil {
		log.Fatal("Database error:", err)
	}
	redisClient, err := redis.Init(config.Redis)
	if err != nil {
		log.Fatal("redis connection error:", err)
	}

	// Create a new router
	router := mux.NewRouter()
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // The URL to fetch the Swagger JSON file
	))
	// Create a private router for authenticated routes
	privateRouter := router.PathPrefix("/").Subrouter()

	// Create instances of your handlers and services
	userRepo := domain.NewUserRepositoryDb(dbClient, redisClient, config.DB.Database, config.DB.UserCollection)
	userService := service.NewUserService(userRepo)
	userHandler := handler.UserHandler{Service: userService}

	// Define routes and corresponding handler methods
	router.HandleFunc("/create-user", userHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/user/sign-in", userHandler.SignIn).Methods(http.MethodPost)
	privateRouter.HandleFunc("/user/get-user", userHandler.GetUser).Methods(http.MethodGet)

	// Setup middleware
	privateRouter.Use(middleware.Authentication)

	// CORS middleware
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"*"})

	// Start the HTTP server
	address := ":8080" // Replace with your desired port
	log.Printf("Server is listening on %s...\n", address)

	// Wrap the router with CORS middleware and start the server
	err = http.ListenAndServe(address, handlers.CORS(originsOk, headersOk, methodsOk)(router))
	if err != nil {
		log.Fatal("HTTP server error:", err)
	}
}
