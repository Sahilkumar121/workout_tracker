package service

import (
	"context"

	"github.com/Sahilkumar121/workout_tracker/internal/handler/user"
	"github.com/Sahilkumar121/workout_tracker/internal/security"
)

type UserRepository interface {
	FindUserByEmail(
		ctx context.Context,
		email string,
	) (*user.UserCredentials, error)

	CreateUser(
		ctx context.Context,
		req *user.RegisterRequest,
	) (*user.UserResponse, error)

	FindUserByID(
		ctx context.Context,
		id string,
	) (*user.UserResponse, error)
}

type UserService struct {
	Repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (s *UserService) CheckValidUserService(
	ctx context.Context,
	email string,
) (*user.UserCredentials, error) {

	d, err := s.Repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *UserService) RegisterUserService(
	ctx context.Context,
	req *user.RegisterRequest,
) (*user.UserResponse, error) {

	hashPass, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	req.Password = hashPass

	d, err := s.Repo.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (s *UserService) FindUserByID(
	ctx context.Context,
	id string,
) (*user.UserResponse, error) {

	d, err := s.Repo.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return d, nil
}
