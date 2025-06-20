package main

import (
	"client/common"
	"client/internal/handler"
	"client/internal/router"
	"log"
	"net/http"
	"os"

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

	client := service.NewServiceJSONClient(os.Getenv("SERVICE_ENDPOINT"), &http.Client{})
	handlers := handler.NewHandler(client)
	router := router.NewRouter(handlers)
	router.Register(e)

	// logger
	e.Use(common.LogRequest)
	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":8083"))
}
