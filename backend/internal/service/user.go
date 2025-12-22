package service

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo repository.UserRepositoryImp
}

func NewUserService(userRepo repository.UserRepositoryImp) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetMe(ctx context.Context, userID uuid.UUID) (*user.UserResponse, error) {
	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &user.UserResponse{
		ID:              u.ID,
		Email:           u.Email,
		Name:            u.Name,
		Role:            u.Role,
		ProfileImageURL: u.ProfileImageURL,
	}, nil
}
