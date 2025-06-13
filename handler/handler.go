package handler

import "github.com/Anhbman/microservice-server-cake/rpc/service"

type Handler struct {
	serviceClient service.Service
}

func NewHandler(serviceClient service.Service) *Handler {
	return &Handler{
		serviceClient: serviceClient,
	}
}
