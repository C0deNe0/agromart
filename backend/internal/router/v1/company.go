package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterCompanyRoutes(r *echo.Group, h *handler.Handlers, auth *middleware.AuthMiddleware) {
	company := r.Group("/companies", auth.RequireAuth())
	company.POST("", h.Company.Create())
	company.GET("/me", h.Company.ListMine())
}
