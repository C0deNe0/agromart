package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/model/auth"
	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type AuthService struct {
	userRepo         *repository.UserRepository
	authMethodRepo   *repository.UserAuthMethodRepository
	TokenManager     *utils.TokenManager
	refreshTokenRepo *repository.RefreshTokenRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
	authMethodRepo *repository.UserAuthMethodRepository,
	tokenManager *utils.TokenManager,
	refreshTokenRepo *repository.RefreshTokenRepository,
) *AuthService {
	return &AuthService{
		userRepo:         userRepo,
		authMethodRepo:   authMethodRepo,
		TokenManager:     tokenManager,
		refreshTokenRepo: refreshTokenRepo,
	}
}

type ctxKey string

const (
	AuthProviderLocal  = "LOCAL"
	AuthProviderGoogle = "GOOGLE"
	ctxUserAgent       = ctxKey("user_agent")
	ctxIPAddress       = ctxKey("ip")
)

// register with email
func (s *AuthService) RegisterWithEmail(ctx context.Context, email string, password string, name string) (*user.AuthResponse, error) {
	// If DB is down → existing == nil → duplicate user creation attempt.
	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, errors.New("email already exists")
	}
	// if err != nil && !errors.Is(err, repository.ErrNotFound) {
	// 	return nil, err
	// }

	u := &user.User{
		Email:         email,
		Name:          name,
		Role:          user.RoleUser,
		IsActive:      true,
		EmailVerified: false,
	}

	createdUser, err := s.userRepo.Create(ctx, u)
	if err != nil {
		return nil, err
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	method := &repository.UserAuthMethod{
		UserId:       createdUser.ID,
		AuthProvider: AuthProviderLocal,
		PasswordHash: &hash,
	}

	_, err = s.authMethodRepo.Create(ctx, method)
	if err != nil {
		return nil, err
	}

	//Issue token
	accessToken, refreshToken, err := s.issueTokens(ctx, createdUser.ID, string(createdUser.Role))
	if err != nil {
		return nil, err
	}

	return &user.AuthResponse{
		User: user.UserResponse{
			ID:    createdUser.ID,
			Email: createdUser.Email,
			Name:  createdUser.Name,
			Role:  createdUser.Role,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// login with  email passowrd
func (s *AuthService) LoginWithEmail(ctx context.Context, email string, password string) (*user.AuthResponse, error) {
	//Get user by email
	method, err := s.authMethodRepo.GetLocalByEmail(ctx, email)
	if err != nil || method.PasswordHash == nil {
		return nil, fmt.Errorf("getlocal email failed %v", err)
	}
	if err := utils.VerifyPassword(*method.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("verify password failed %v", err)
	}

	u, err := s.userRepo.GetByID(ctx, method.UserId)
	if err != nil || !u.IsActive {
		return nil, fmt.Errorf("user not allowed %v", err)
	}

	//Issue token
	accessToken, refreshToken, err := s.issueTokens(ctx, u.ID, string(u.Role))
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

func (s *AuthService) LoginWithGoogle(ctx context.Context, googleSub string, email, name string, profileURL *string) (*user.AuthResponse, error) {
	//Get user by email
	u, err := s.userRepo.GetByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		u, err = s.userRepo.Create(ctx, &user.User{
			Email:           email,
			Name:            name,
			Role:            user.RoleUser,
			EmailVerified:   true,
			IsActive:        true,
			ProfileImageURL: profileURL,
		})
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if _, err := s.authMethodRepo.EnsureOAuth(ctx, u.ID, AuthProviderGoogle, googleSub); err != nil {
		return nil, err
	}

	access, refresh, err := s.issueTokens(ctx, u.ID, string(u.Role))
	if err != nil {
		return nil, err
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

func (s *AuthService) Refresh(ctx context.Context, rawRefreshToken string) (*user.AuthResponse, error) {

	// 1. Hash incoming refresh token
	tokenHash := utils.HashToken(rawRefreshToken)

	// 2. Find valid refresh token in DB
	rt, err := s.refreshTokenRepo.FindValid(ctx, tokenHash)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	if rt.RevokedAt != nil {
		return nil, errors.New("refresh token revoked")
	}

	// 3. Load user
	u, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil || !u.IsActive {
		return nil, errors.New("user not allowed")
	}

	// 4. Revoke old refresh token (rotation)
	if err := s.refreshTokenRepo.Revoke(ctx, tokenHash); err != nil {
		return nil, err
	}

	// 5. Issue new tokens + persist new refresh token
	access, refresh, err := s.issueTokens(
		ctx,
		u.ID,
		string(u.Role),
	)
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
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, rawRefreshToken string) error {

	refreshToken := utils.HashToken(rawRefreshToken)

	if err := s.refreshTokenRepo.Revoke(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}

// this is just a auth helper function
func (s *AuthService) issueTokens(ctx context.Context, userID uuid.UUID, role string) (access string, refresh string, err error) {

	access, err = s.TokenManager.GenerateAccessToken(userID, role)
	if err != nil {
		return
	}

	refresh, err = s.TokenManager.GenerateRefreshToken(userID)
	if err != nil {
		return
	}

	_, err = s.refreshTokenRepo.Create(ctx, &auth.RefreshToken{
		UserID:    userID,
		TokenHash: utils.HashToken(refresh),
		UserAgent: getCtxString(ctx, "user_agent"),
		IPAddress: getCtxString(ctx, "ip"),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	})

	return
}

func getCtxString(ctx context.Context, key ctxKey) string {
	if v := ctx.Value(key); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
