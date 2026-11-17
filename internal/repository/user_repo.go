package repository

import (
	"context"

	"github.com/Sahilkumar121/workout_tracker/internal/handler/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) FindUserByEmail(
	ctx context.Context,
	email string,
) (*user.UserCredentials, error) {

	query := `select id, password from users where email = $1`

	var res user.UserCredentials
	if err := r.DB.QueryRow(ctx, query, email).Scan(
		&res.ID,
		&res.Password,
	); err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *UserRepository) FindUserByID(
	ctx context.Context,
	id string,
) (*user.UserResponse, error) {

	query := `select id, name, email, created_at, updated_at from users where id = $1`

	var userRes user.UserResponse
	if err := r.DB.QueryRow(ctx, query, id).Scan(
		&userRes.ID,
		&userRes.Name,
		&userRes.Email,
		&userRes.CreatedAt,
		&userRes.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &userRes, nil
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	req *user.RegisterRequest,
) (*user.UserResponse, error) {

	query := `
			insert into users(name, email, password)
			values ($1, $2, $3)
			returning id, name, email, created_at, updated_at;
			`

	var userRes user.UserResponse
	err := r.DB.QueryRow(ctx, query, req.Name, req.Email, req.Password).Scan(
		&userRes.ID,
		&userRes.Name,
		&userRes.Email,
		&userRes.CreatedAt,
		&userRes.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &userRes, nil
}
