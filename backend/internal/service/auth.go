package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo       repository.UserRepositoryImp
	authMethodRepo repository.UserAuthMethodRepositoryImp
	TokenManager   *utils.TokenManager
}

func NewAuthService(userRepo repository.UserRepositoryImp, authMethodRepo repository.UserAuthMethodRepositoryImp, tokenManager *utils.TokenManager) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		authMethodRepo: authMethodRepo,
		TokenManager:   tokenManager,
	}
}

// register with email
func (s *AuthService) RegisterWithEmail(ctx context.Context, email string, password string, name string) (*user.AuthResponse, error) {
	existing, _ := s.userRepo.GetByEmail(ctx, email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	u := &user.User{
		Email:    email,
		Name:     name,
		Role:     user.UserRole(user.RoleUser),
		IsActive: true,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, err
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	method := &repository.UserAuthMethod{
		ID:           uuid.New(),
		UserId:       u.ID,
		Provider:     "EMAIL_PASSWORD",
		PasswordHash: &hash,
	}

	if err := s.authMethodRepo.Create(ctx, method); err != nil {

		return nil, err
	}

	//Issue token
	accessToken, err := s.TokenManager.GenerateAccessToken(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.TokenManager.GenerateRefreshToken(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	return &user.AuthResponse{
		User: user.UserResponse{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.Name,
			Role:  u.Role,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// login with  email passowrd
func (s *AuthService) LoginWithEmail(ctx context.Context, email string, password string) (*user.AuthResponse, error) {
	method, err := s.authMethodRepo.GetByEmail(ctx, email)
	if err != nil || method.PasswordHash == nil {
		return nil, errors.New("invalid credentials")
	}
	if err := utils.VerifyPassword(*method.PasswordHash, password); err != nil {
		return nil, errors.New("Invalid credentials")
	}

	u, err := s.userRepo.GetByID(ctx, method.UserId)
	if err != nil {
		return nil, err
	}

	//Issue token
	accessToken, err := s.TokenManager.GenerateAccessToken(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.TokenManager.GenerateRefreshToken(u.ID, string(u.Role))
	if err != nil {
		return nil, err
	}

	return &user.AuthResponse{
		User: user.UserResponse{
			ID:    u.ID,
			Email: u.Email,
			Name:  u.Name,
			Role:  u.Role,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) LoginWithGoogle(ctx context.Context, googleUserID string, email, name string, profileURL *string) (*user.AuthResponse, error) {
	u, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || u == nil {
		u = &user.User{
			Email:    email,
			Name:     name,
			Role:     user.RoleUser,
			IsActive: true,
		}

		if err := s.userRepo.Create(ctx, u); err != nil {
			return nil, err
		}
	}

	if err := s.authMethodRepo.EnsureOAuth(ctx, u.ID, "GOOGLE", googleUserID); err != nil {
		return nil, err
	}

	access, err := s.TokenManager.GenerateAccessToken(u.ID, string(u.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %v", err)
	}

	refresh, err := s.TokenManager.GenerateRefreshToken(u.ID, string(u.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return &user.AuthResponse{
		User: user.UserResponse{
			ID:              u.ID,
			Email:           u.Email,
			Name:            u.Name,
			Role:            u.Role,
			ProfileImageURL: u.ProfileImageURL,
		},
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}
