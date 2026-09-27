package handler

import (
	"time"
)

// user
type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterUser struct {
	FirstName string `json:"first_name" validate:"required,min=2,alpha,max=50"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,min=3,alpha,max=50"`
	Email     string `json:"email" validate:"required,email,max=100"`
	Password  string `json:"password" validate:"required,min=12,max=100"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=12,max=100"`
}

type LoginUserResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
}

type UpdateUserData struct {
	FirstName string `json:"first_name,omitempty" validate:"omitempty,alpha,max=50"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,alpha,.max=50"`
	Email     string `json:"email,omitempty" validate:"omitempty,email,max=100"`
	Password  string `json:"password,omitempty" validate:"omitempty,min=12,max=100"`
}

// health
type HealthResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

// workout
type WorkoutRequest struct {
	WorkoutName      string     `json:"workout_name" validate:"required,max=100"`
	WorkoutStartTime *time.Time `json:"workout_start,omitempty" validate:"omitempty"`
}

type WorkoutResponse struct {
	ID           string     `json:"id"`
	WorkoutName  string     `json:"workout_name"`
	WorkoutStart *time.Time `json:"workout_start"`
}

type WorkoutRequestByName struct {
	WorkoutName string `json:"workout_name" validate:"required"`
}
