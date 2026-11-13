package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Sahilkumar121/workout_tracker/internal/handler/exercise"
)

type ExerciseRepository interface {
	InsertIntoExercise(
		ctx context.Context,
		req exercise.PostRequest,
	) (*exercise.Exercise, error)
}

type ExerciseServiceHandler struct {
	Logger *slog.Logger
	Repo   ExerciseRepository
}

func NewExerciseServiceHandler(logger *slog.Logger, repo ExerciseRepository) *ExerciseServiceHandler {
	return &ExerciseServiceHandler{
		Logger: logger,
		Repo:   repo,
	}
}

func (s *ExerciseServiceHandler) CreateExercise(
	ctx context.Context,
	exerciseBodyReq exercise.PostRequest,
) (*exercise.Exercise, error) {

	if strings.TrimSpace(exerciseBodyReq.Name) == "" ||
		strings.TrimSpace(exerciseBodyReq.Description) == "" ||
		strings.TrimSpace(exerciseBodyReq.MuscleGroup) == "" {
		return nil, fmt.Errorf("name and description was not provide")
	}

	response, err := s.Repo.InsertIntoExercise(ctx, exerciseBodyReq)
	if err != nil {
		return nil, fmt.Errorf("create exercise error: %v\n", err)
	}

	return response, nil
}
