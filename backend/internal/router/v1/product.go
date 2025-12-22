package v1

import (
	"github.com/C0deNe0/agromart/internal/handlers"
	"github.com/labstack/echo/v4"
)

func RegisterProductRoutes(r *echo.Group, h *handlers.Handlers) {
	product := r.Group("/products")
	// product.GET("/", h.Product.ListProducts(echo.Context))
}
