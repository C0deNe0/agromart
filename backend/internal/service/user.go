package service

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetMe(ctx context.Context, userID uuid.UUID) (*user.UserResponse, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := user.UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		Name:            u.Name,
		Role:            u.Role,
		ProfileImageURL: u.ProfileImageURL,
	}

	return &resp, nil
}

func (s *UserService) BlockUser(ctx context.Context, userID uuid.UUID) error {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	u.IsActive = false
	_, err = s.userRepo.Update(ctx, u)
	return err
}
