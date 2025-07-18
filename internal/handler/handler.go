package handler

import (
	"client/internal/service"

	serverService "github.com/Anhbman/microservice-server-cake/rpc/service"
)

type Handler struct {
	serviceClient serverService.Service
	eventService  *service.EventService
}

func NewHandler(serviceClient serverService.Service, eventService *service.EventService) *Handler {
	return &Handler{
		serviceClient: serviceClient,
		eventService:  eventService,
	}
}
