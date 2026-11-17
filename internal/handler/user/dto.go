package user

import (
	"time"

	"github.com/google/uuid"
)

// Incoming register request
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,alphaspace,max=100,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Incoming login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Outgoing user response
type UserResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// Internal data required for authentication
type UserCredentials struct {
	ID       uuid.UUID
	Password string
}
