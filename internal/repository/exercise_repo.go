package repository

import (
	"context"
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
) (exercise.Exercise, error) {

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
		return exercise.Exercise{}, fmt.Errorf("insertion into exercise failed %w", err)
	}

	r.Logger.Info("insertion into exercise successfully")
	return response, nil
}

func (r *ExerciseRepository) FindAllExercises(
	ctx context.Context,
) ([]exercise.Exercise, error) {

	query := `select id, name, description, muscle_group, created_at from exercises;`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("database query: %w", err)
	}

	defer rows.Close()

	d, err := pgx.CollectRows(rows, pgx.RowToStructByName[exercise.Exercise])
	if err != nil {
		return nil, fmt.Errorf("collect row: %vw", err)
	}

	return d, nil
}

func (r *ExerciseRepository) FindExerciseById(
	ctx context.Context,
	id string,
) (exercise.Exercise, error) {

	query := `select id, name, description, muscle_group, created_at from exercises where id = $1;`

	var response exercise.Exercise
	if err := r.DB.QueryRow(ctx, query, id).Scan(
		&response.Id,
		&response.Name,
		&response.Description,
		&response.MuscleGroup,
		&response.CreatedAt,
	); err != nil {

		return exercise.Exercise{}, fmt.Errorf("find exercise by id: %w", err)
	}

	return response, nil
}
