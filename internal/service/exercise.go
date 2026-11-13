package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Sahilkumar121/workout_tracker/internal/handler/exercise"
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

	UpdateExerciseById(
		ctx context.Context,
		id string,
		request exercise.PatchRequest,
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

	d, err := s.Repo.FindExerciseById(ctx, id)

	if err != nil {
		return exercise.Exercise{}, err
	}

	return d, nil
}

func (s *ExerciseServiceHandler) UpdateExerciseById(
	ctx context.Context,
	id string,
	patchReq exercise.PatchRequest,
) (exercise.Exercise, error) {

	d, err := s.Repo.UpdateExerciseById(ctx, id, patchReq)

	if err != nil {
		return exercise.Exercise{}, err
	}

	return d, nil
}
