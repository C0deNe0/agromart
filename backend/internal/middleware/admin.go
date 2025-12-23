package middleware

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/labstack/echo/v4"
)

type AdminMiddleware struct {
}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (a *AdminMiddleware) RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role := GetUserRole(c)
		if role != user.RoleAdmin {
			return echo.NewHTTPError(http.StatusForbidden, "You are not authorized to access this resource")
		}
		return next(c)
	}
}
