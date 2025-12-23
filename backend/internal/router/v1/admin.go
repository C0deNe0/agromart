package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterAdminRoutes(r *echo.Group, h *handler.Handlers) {

	r.GET("/companies/pending", h.Company.ListPending)
	r.POST("/companies/:id/approve", h.Company.Approve)
	r.POST("/companies/:id/reject", h.Company.Reject)
}
