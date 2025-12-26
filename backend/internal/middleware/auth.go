package middleware

import (
	"fmt"
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

			tokenStr := extractBearerToken(c)
			if tokenStr == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}
			claims, err := am.tokenManager.ParseAccessToken(tokenStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("Unauthorized: %v", err))
			}
			c.Set("userID", claims.UserID)
			c.Set("role", user.UserRole(claims.Role))
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

	token := authHeader[len(bearerPrefix):]
	token = strings.TrimSpace(token)
	token = strings.Trim(token, "\"")
	return token
}


// . The Server Crash (Panic)
// The initial "access token problem" you reported was causing your Go server to crash with a panic whenever a protected route was hit. This happened because of a mismatch between your Middleware and your Helper functions:

// Wrong Key Name: Your middleware was saving the user ID as "user_id" (c.Set("user_id", ...)), but your helper function middleware.GetUserID(c) was looking for "userID". This meant the helper found nothing (nil).
// Wrong Data Type: Your middleware saved the role as a plain string. However, your helper middleware.GetUserRole(c) tried to force-convert it into a user.UserRole (a custom type). In Go, trying to convert an interface holding string directly to user.UserRole without an explicit cast causes a panic (crash).
// The Fix: I updated 
// auth.go
//  to use the correct keys ("userID") and store the role as the correct type (user.UserRole). I also updated 
// context.go
//  to be "panic-safe," meaning if data is missing, it now just returns an empty value instead of crashing the entire server.

// 2. The Malformed Token (401 Unauthorized)
// After the crash was fixed, you saw an error: illegal base64 data at input byte 0.

// The Cause: This error came from the JWT library. It tried to decode your token but found a character that wasn't valid Base64 right at the beginning (byte 0).
// Why? This almost always happens when you copy-paste a token from a JSON response and accidentally include the surrounding double quotes (e.g., "eyJhb..."). A quote " is not a valid character for a JWT signature.