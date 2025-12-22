package v1

import (
	"github.com/C0deNe0/agromart/internal/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterCompanyRoutes(r *echo.Group, h *handlers.Handlers) {
	company := r.Group("/companies")
	company.GET("/", h.Company.Create)
	company.POST("/", h.Company.ListMine)
}
