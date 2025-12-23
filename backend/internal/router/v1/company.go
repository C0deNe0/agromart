package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterCompanyRoutes(c echo.Context, r *echo.Group, h *handler.Handlers) {
	company := r.Group("/companies")
	company.GET("/", h.Company.Create(c))
	// company.POST("/", h.Company.ListMine(c))
}
