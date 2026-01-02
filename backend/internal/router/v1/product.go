package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterProductRoutes(r *echo.Group, h *handler.Handlers) {

	product := r.Group("/products")

	product.GET("/", h.Product.ListProducts())
	product.GET("/with-category", h.Product.ListProductsWithCategory())
	product.GET("/:id", h.Product.GetProductByID())

	product.POST("/", h.Product.CreateProduct())
	product.PUT("/:id", h.Product.UpdateProduct())

	product.DELETE("/:id", h.Product.DeleteProduct())
}
