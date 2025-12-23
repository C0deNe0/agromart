package router

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterSystemRoutes(r *echo.Echo, h *handler.Handlers) {
	r.GET("/status", h.Health.CheckHealth)
}
