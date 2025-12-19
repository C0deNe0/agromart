package user

import (
	"time"

	"github.com/C0deNe0/agromart/internal/model"
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
