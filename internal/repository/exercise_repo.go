package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Sahilkumar121/workout_tracker/internal/handler/exercise"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExerciseRepository struct {
	Logger *slog.Logger
	DB     *pgxpool.Pool
}

func NewExerciseRepository(logger *slog.Logger, db *pgxpool.Pool) *ExerciseRepository {
	return &ExerciseRepository{
		Logger: logger,
		DB:     db,
	}
}

func (r *ExerciseRepository) InsertIntoExercise(
	ctx context.Context,
	req exercise.PostRequest,
) (*exercise.Exercise, error) {

	query := `
		insert into exercises(name, description, muscle_group)
		values ($1, $2, $3) 
		returning id, name, description, muscle_group, created_at;
	`

	var response exercise.Exercise
	err := r.DB.QueryRow(ctx, query, req.Name, req.Description, req.MuscleGroup).Scan(
		&response.Id,
		&response.Name,
		&response.Description,
		&response.MuscleGroup,
		&response.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no row found %v\n", err)
		}
		return nil, fmt.Errorf("insertion into exercise failed %v\n", err)
	}

	r.Logger.Info("insertion into exercise successfully")
	return &response, nil
}
