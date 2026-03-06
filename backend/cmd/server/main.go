package main

import (
	"log"
	"net/http"
	"os"
	"take-Home-assignment/internal/config"
	"take-Home-assignment/internal/database"
	"take-Home-assignment/internal/handlers"
	"take-Home-assignment/internal/routes"
	"take-Home-assignment/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.LoadConfig()

	database.Connect()

	r := gin.Default()

	serviceContainer := services.InitServices(database.DB, cfg)

	handlerContainer := handlers.InitHandlers(serviceContainer)

	routes.Routes(r, handlerContainer)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	server.ListenAndServe()
}
