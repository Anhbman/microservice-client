package messaging

const (
	Topic = "topic"
)

// Define event types as constants
const (
	UserRegisteredEvent = "user.registered"
	UserUpdatedEvent    = "user.updated"
	UserDeletedEvent    = "user.deleted"
	CakeCreatedEvent    = "cake.created"
	CakeUpdatedEvent    = "cake.updated"
	CakeDeletedEvent    = "cake.deleted"
	OrderCreatedEvent   = "order.created"
	OrderUpdatedEvent   = "order.updated"
)

// Event payloads
type UserRegisteredPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdatedPayload struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CakeCreatedPayload struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	UserID      uint64 `json:"user_id"`
}
