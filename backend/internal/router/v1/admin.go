package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterAdminRoutes(r *echo.Group, h *handler.Handlers, admin *middleware.AdminMiddleware) {

	adminGroup := r.Group("/admin")
	adminGroup.Use(admin.RequireAdmin)

	adminGroup.PUT("/companies/:id/approve", h.Admin.ApproveCompany())
	adminGroup.PUT("/companies/:id/reject", h.Admin.RejectCompany())
	adminGroup.GET("/companies/pending", h.Admin.CountPendingApprovals())
	// adminGroup.DELETE("/companies/:id", h.Admin.DeleteCompany())
}
