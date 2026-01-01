package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterProductRoutes(r *echo.Group, h *handler.Handlers) {
	product := r.Group("/products")
	product.GET("/", h.Product.ListProducts())
	product.POST("/", h.Product.CreateProduct())
	product.GET("/:id", h.Product.GetProductByID())
	product.PUT("/:id", h.Product.UpdateProduct())
	product.DELETE("/:id", h.Product.DeleteProduct())

	// product.POST("/:id/images", h.Product.GenerateImageUploadURL())
	// product.POST("/:id/images/presign", h.Product.GenerateImageUploadURL())
	// product.DELETE("/:id/images/:imageID", h.Product.DeleteImage())
}
