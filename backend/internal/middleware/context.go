package middleware

import "github.com/google/uuid"

func GetUserID(c interface {
	Get(string) interface{}
}) uuid.UUID {
	return c.Get("userID").(uuid.UUID)
}

func GetUserRole(c interface {
	Get(string) interface{}
}) string {
	return c.Get("role").(string)
}

func IsAdmin(c interface {
	Get(string) interface{}
}) bool {
	return c.Get("role").(string) == "ADMIN"
}
