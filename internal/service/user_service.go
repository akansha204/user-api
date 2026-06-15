package service

import (
	"context"
	"errors"
	"time"

	"user-api/db/sqlc"
	"user-api/internal/models"
	"user-api/internal/repository"
)

const dobLayout = "2006-01-02"
const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
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
	if s.repo == nil {
		return models.UserResponse{}, errors.New("repository not configured")
	}

	dob, err := time.Parse(dobLayout, req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.CreateUser(ctx, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, false, s.now()), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (models.UserResponse, error) {
	if s.repo == nil {
		return models.UserResponse{}, errors.New("repository not configured")
	}

	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, true, s.now()), nil
}

func (s *UserService) ListUsers(ctx context.Context, page, limit int) ([]models.UserResponse, error) {
	if s.repo == nil {
		return nil, errors.New("repository not configured")
	}

	page, limit = normalizePagination(page, limit)
	offset := int32((page - 1) * limit)

	users, err := s.repo.ListUsers(ctx, int32(limit), offset)
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
	if s.repo == nil {
		return models.UserResponse{}, errors.New("repository not configured")
	}

	dob, err := time.Parse(dobLayout, req.Dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := s.repo.UpdateUser(ctx, id, req.Name, dob)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, false, s.now()), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	if s.repo == nil {
		return errors.New("repository not configured")
	}

	return s.repo.DeleteUser(ctx, id)
}

func normalizePagination(page, limit int) (int, int) {
	if page < 1 {
		page = defaultPage
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return page, limit
}

func toUserResponse(user sqlc.User, includeAge bool, now time.Time) models.UserResponse {
	response := models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Dob:  user.Dob.UTC().Format(dobLayout),
	}

	if includeAge {
		response.Age = CalculateAge(user.Dob, now)
	}

	return response
}
