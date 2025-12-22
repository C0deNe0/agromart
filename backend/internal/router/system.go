package router

import (
	"github.com/C0deNe0/agromart/internal/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterSystemRoutes(r *echo.Echo, h *handlers.Handlers) {
	r.GET("/status", h.Health.CheckHealth)
}
