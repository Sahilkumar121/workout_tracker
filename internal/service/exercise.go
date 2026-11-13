package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Sahilkumar121/workout_tracker/internal/handler/exercise"
	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/google/uuid"
)

type ExerciseRepository interface {
	InsertIntoExercise(
		ctx context.Context,
		req exercise.PostRequest,
	) (exercise.Exercise, error)

	FindAllExercises(
		ctx context.Context,
	) ([]exercise.Exercise, error)

	FindExerciseById(
		ctx context.Context,
		id string,
	) (exercise.Exercise, error)
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
) (exercise.Exercise, error) {

	if strings.TrimSpace(exerciseBodyReq.Name) == "" ||
		strings.TrimSpace(exerciseBodyReq.Description) == "" ||
		strings.TrimSpace(exerciseBodyReq.MuscleGroup) == "" {
		return exercise.Exercise{}, httpx.ErrBadRequest
	}

	response, err := s.Repo.InsertIntoExercise(ctx, exerciseBodyReq)
	if err != nil {
		return exercise.Exercise{}, fmt.Errorf("create exercise error: %w", err)
	}

	return response, nil
}

func (s *ExerciseServiceHandler) FindAllExercise(
	ctx context.Context,
) ([]exercise.Exercise, error) {

	d, err := s.Repo.FindAllExercises(ctx)
	if err != nil {
		return []exercise.Exercise{}, fmt.Errorf("find all exercise failed: %v\n", err)
	}

	return d, nil
}

func (s *ExerciseServiceHandler) SearchExerciseByID(
	ctx context.Context,
	id string,
) (exercise.Exercise, error) {

	if strings.TrimSpace(id) == "" {
		return exercise.Exercise{}, fmt.Errorf("id can't be empty")
	}

	if _, err := uuid.Parse(id); err != nil {
		return exercise.Exercise{}, fmt.Errorf("invalid id: %w", err)
	}

	d, err := s.Repo.FindExerciseById(ctx, id)
	if err != nil {
		return exercise.Exercise{}, fmt.Errorf("search exercise by id: %w", err)
	}

	return d, nil
}
