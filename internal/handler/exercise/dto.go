package exercise

import (
	"time"

	"github.com/google/uuid"
)

type Exercise struct {
	Id          uuid.UUID  `json:"id" validate:"required,uuid"`
	Name        string     `json:"name" validate:"required,alphaspace"`
	Description string     `json:"description" validate:"required"`
	MuscleGroup string     `json:"muscle_group" validate:"required"`
	CreatedAt   *time.Time `json:"created_at" validate:"required"`
}

type PostRequest struct {
	Name        string `json:"name" validate:"required,alphaspace"`
	Description string `json:"description" validate:"required"`
	MuscleGroup string `json:"muscle_group" validate:"required"`
}

type PatchRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	MuscleGroup *string `json:"muscle_group"`
}
