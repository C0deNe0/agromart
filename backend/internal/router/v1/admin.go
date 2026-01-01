package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterAdminRoutes(r *echo.Group, h *handler.Handlers, admin *middleware.AdminMiddleware) {

	adminGroup := r.Group("/admin")
	adminGroup.Use(admin.RequireAdmin)

	adminGroup.DELETE("/companies/:id", h.Company.Delete())
	adminGroup.GET("/companies/pending", h.Company.ListPending())
	adminGroup.POST("/companies/:id/approve", h.Company.Approve())
	adminGroup.POST("/companies/:id/reject", h.Company.Reject())
}
