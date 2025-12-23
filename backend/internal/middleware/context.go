package middleware

import (
	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/google/uuid"
)

func GetUserID(c interface {
	Get(string) interface{}
}) uuid.UUID {
	return c.Get("userID").(uuid.UUID)
}

func GetUserRole(c interface {
	Get(string) interface{}
}) user.UserRole {
	return c.Get("role").(user.UserRole)
}

func IsAdmin(c interface {
	Get(string) interface{}
}) bool {
	return c.Get("role").(user.UserRole) == user.RoleAdmin
}
