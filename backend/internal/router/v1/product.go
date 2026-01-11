package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterProductRoutes(r *echo.Group, h *handler.Handlers) {

	product := r.Group("/products")

	product.GET("/", h.Product.ListProducts())
	product.GET("/:id", h.Product.GetProductByID())

	product.POST("/", h.Product.CreateProduct())
	product.PUT("/:id", h.Product.UpdateProduct())
	product.DELETE("/:id", h.Product.DeleteProduct())
	product.POST("/:id/resubmit", h.Product.ResubmitProduct())
	product.GET("/:id/history", h.Product.GetApprovalHistory())

	//IMage
	product.POST("/:id/images/upload-url", h.Product.GenerateImageUploadURL())
	product.DELETE("/:id/images/:imageId", h.Product.DeleteImage())
	product.PUT("/:id/images/:imageId/primary", h.Product.SetPrimaryImage())

	//variant
	product.POST("/:id/variants", h.Product.CreateVariant())
	product.PUT("/:id/variants/:variantId", h.Product.UpdateVariant())
	product.DELETE("/:id/variants/:variantId", h.Product.DeleteVariant())

}
