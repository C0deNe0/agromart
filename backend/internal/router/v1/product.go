package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterProductRoutes(c echo.Context, r *echo.Group, h *handler.Handlers) {
	product := r.Group("/products")
	product.GET("/", h.Product.ListProducts(c))
}
