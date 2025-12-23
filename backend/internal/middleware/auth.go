package middleware

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/labstack/echo/v4"
)

type AuthMiddleware struct {
	tokenManager *utils.TokenManager
}

func NewAuthMiddleware(tokenManager *utils.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{tokenManager: tokenManager}
}

func (am *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Implement your authentication logic here

			tokenStr := extractBearerToken(c)
			if tokenStr == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			claims, err := am.tokenManager.ParseAccessToken(tokenStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)
			return next(c)
		}
	}
}

func extractBearerToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return ""
	}

	return authHeader[len(bearerPrefix):]
}
