package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminOnly(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !IsAdmin(c) {
			return echo.NewHTTPError(http.StatusForbidden, "You are not authorized to access this resource")
		}
		return next(c)
	}
}
