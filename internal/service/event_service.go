package service

import (
	"client/internal/messaging"
	"context"

	"github.com/Anhbman/microservice-server-cake/rpc/service"
)

type EventService struct {
	eventManager *messaging.EventManager
}

func NewEventService(eventManager *messaging.EventManager) *EventService {
	return &EventService{
		eventManager: eventManager,
	}
}

func (s *EventService) PublishUserRegistered(ctx context.Context, req *service.RegisterUserRequest) error {
	payload := messaging.UserRegisteredPayload{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	return s.eventManager.PublishEvent(ctx, messaging.UserRegisteredEvent, payload)
}
