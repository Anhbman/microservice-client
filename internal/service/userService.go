package service

import (
	"client/internal/messaging"

	"github.com/Anhbman/microservice-server-cake/rpc/service"
)

type UserService struct {
	publisher *messaging.Publisher
}

func NewUserService(publisher *messaging.Publisher) *UserService {
	return &UserService{
		publisher: publisher,
	}
}

func (s *UserService) RegisterUser(req *service.RegisterUserRequest) error {
	// Create payload for RabbitMQ
	payload := messaging.UserRegisterPayload{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// Publish event to RabbitMQ
	return s.publisher.PublishUserRegisteredEvent(payload)
}
