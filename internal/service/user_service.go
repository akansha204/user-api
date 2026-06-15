package service

import (
	"context"
	"time"

	"user-api/db/sqlc"
	"user-api/internal/models"
	"user-api/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
	now  func() time.Time
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
		now:  time.Now,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserResponse, error) {
	user, err := s.repo.CreateUser(ctx, req.Name, req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, false, s.now()), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (models.UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, true, s.now()), nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, toUserResponse(user, true, s.now()))
	}

	return responses, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserResponse, error) {
	user, err := s.repo.UpdateUser(ctx, id, req.Name, req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, false, s.now()), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.repo.DeleteUser(ctx, id)
}

func toUserResponse(user sqlc.User, includeAge bool, now time.Time) models.UserResponse {
	response := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob,
	}

	if includeAge {
		response.Age = CalculateAge(user.Dob, now)
	}

	return response
}
