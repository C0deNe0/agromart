package user

import "github.com/C0deNe0/agromart/internal/model"

type UserRole string

const (
	RoleUser  UserRole = "USER"
	RoleAdmin UserRole = "ADMIN"
)

type User struct {
	model.Base
	Email string `json:"email" db:"email"`
	Name  string `json:"name" db:"name"`
	
}
