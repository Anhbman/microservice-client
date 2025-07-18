package main

import (
	"client/internal/common"
	"client/internal/config"
	"client/internal/handler"
	"client/internal/messaging"
	"client/internal/router"
	"log"
	"net/http"
	"os"

	clientService "client/internal/service"

	"github.com/Anhbman/microservice-server-cake/rpc/service"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	e := echo.New()

	// Initialize RabbitMQ
	rabbitMQConfig, err := config.LoadRabbitMQConfig()
	if err != nil {
		log.Fatalf("Failed to load RabbitMQ config: %v", err)
	}

	eventManager, err := messaging.NewEventManager(*rabbitMQConfig)
	if err != nil {
		log.Fatalf("Failed to create event manager: %v", err)
	}

	// Register all events
	if err := eventManager.RegisterEvents(); err != nil {
		log.Fatalf("Failed to register events: %v", err)
	}

	eventService := clientService.NewEventService(eventManager)

	client := service.NewServiceJSONClient(os.Getenv("SERVICE_ENDPOINT"), &http.Client{})
	handlers := handler.NewHandler(client, eventService)
	router := router.NewRouter(handlers)
	router.Register(e)

	// logger
	e.Use(common.LogRequest)
	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":8083"))
}
