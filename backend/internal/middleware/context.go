package middleware

import (
	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/google/uuid"
)

func GetUserID(c interface {
	Get(string) interface{}
}) uuid.UUID {
	val := c.Get("userID")
	if val == nil {
		return uuid.Nil
	}
	return val.(uuid.UUID)
}

func GetUserRole(c interface {
	Get(string) interface{}
}) user.UserRole {
	val := c.Get("role")
	if val == nil {
		return ""
	}
	return val.(user.UserRole)
}

func IsAdmin(c interface {
	Get(string) interface{}
}) bool {
	val := c.Get("role")
	if val == nil {
		return false
	}
	return val.(user.UserRole) == user.RoleAdmin
}
