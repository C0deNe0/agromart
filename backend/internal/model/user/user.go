package user

import (
	"time"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

type User struct {
	model.Base
	Email           string     `json:"email" db:"email"`
	Name            string     `json:"name" db:"name"`
	ProfileImageURL *string    `json:"profileImageURL,omitempty" db:"profile_image_url"`
	Phone           *string    `json:"phone,omitempty" db:"phone"`
	Role            UserRole   `json:"role" db:"role"`
	IsActive        bool       `json:"isActive" db:"is_active"`
	EmailVerified   bool       `json:"emailVerified" db:"email_verified"`
	LastLoginAt     *time.Time `json:"lastLoginAt,omitempty" db:"last_login_at"`
}

type UserResponse struct {
	ID              uuid.UUID `json:"id"`
	Email           string    `json:"email"`
	Name            string    `json:"name"`
	Role            UserRole  `json:"role"`
	ProfileImageURL *string   `json:"profileImageURL,omitempty"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
}
