package main

import (
	"client/internal/common"
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
	rabbitMQConfig := messaging.Config{
		Host:        "localhost",
		Port:        5672,
		Username:    "guest",
		Password:    "guest",
		VirtualHost: "/",
	}

	conn, err := messaging.NewConnection(rabbitMQConfig)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	// Initialize Publisher
	publisher, err := messaging.NewPublisher(conn)
	if err != nil {
		log.Fatalf("Failed to create publisher: %v", err)
	}

	userService := clientService.NewUserService(publisher)

	client := service.NewServiceJSONClient(os.Getenv("SERVICE_ENDPOINT"), &http.Client{})
	handlers := handler.NewHandler(client, userService)
	router := router.NewRouter(handlers)
	router.Register(e)

	// logger
	e.Use(common.LogRequest)
	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":8083"))
}
