package middleware

import (
	"net/http"
	"strings"

	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/model/user"
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

			tokenStr, err := ExtractBearerToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			claims, err := am.tokenManager.ParseAccessToken(tokenStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			c.Set("userID", claims.UserID)
			c.Set("role", user.UserRole(claims.Role))
			return next(c)
		}
	}
}

func ExtractBearerToken(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	token := authHeader[len(bearerPrefix):]
	token = strings.TrimSpace(token)
	token = strings.Trim(token, "\"")
	return token, nil
}
