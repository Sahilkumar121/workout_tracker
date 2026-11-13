package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/Sahilkumar121/workout_tracker/internal/handler/exercise"
	"github.com/Sahilkumar121/workout_tracker/internal/httpx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExerciseRepository struct {
	DB *pgxpool.Pool
}

func NewExerciseRepository(db *pgxpool.Pool) *ExerciseRepository {
	return &ExerciseRepository{
		DB: db,
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

		return exercise.Exercise{}, err
	}

	return response, nil
}

func (r *ExerciseRepository) UpdateExerciseById(
	ctx context.Context,
	id string,
	request exercise.PatchRequest,
) (exercise.Exercise, error) {

	var args []any
	var updates []string

	// check validation of each field
	if request.Name != nil {
		args = append(args, *request.Name)
		updates = append(updates, fmt.Sprintf("name = $%d", len(args)))
	}

	if request.Description != nil {
		args = append(args, *request.Description)
		updates = append(updates, fmt.Sprintf("description = $%d", len(args)))
	}

	if request.MuscleGroup != nil {
		args = append(args, *request.MuscleGroup)
		updates = append(updates, fmt.Sprintf("muscle_group = $%d", len(args)))
	}

	// check if user provided at least one field else return error
	if len(updates) == 0 {
		return exercise.Exercise{}, httpx.ErrNoFieldsToUpdate
	}

	args = append(args, id)

	query := fmt.Sprintf(`
		update exercises
		set %s
		where id = $%d
		returning id, name, description, muscle_group, created_at`,
		strings.Join(updates, ", "), len(args),
	)

	var exerciseRes exercise.Exercise
	if err := r.DB.QueryRow(ctx, query, args...).Scan(
		&exerciseRes.Id,
		&exerciseRes.Name,
		&exerciseRes.Description,
		&exerciseRes.MuscleGroup,
		&exerciseRes.CreatedAt,
	); err != nil {
		return exercise.Exercise{}, err
	}

	return exerciseRes, nil
}
