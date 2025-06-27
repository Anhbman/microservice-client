package handler

import (
	"client/internal/service"

	serverService "github.com/Anhbman/microservice-server-cake/rpc/service"
)

type Handler struct {
	serviceClient serverService.Service
	userService   *service.UserService
}

func NewHandler(serviceClient serverService.Service, userService *service.UserService) *Handler {
	return &Handler{
		serviceClient: serviceClient,
		userService:   userService,
	}
}
