package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterUploadRoutes(r *echo.Group, h *handler.Handlers, auth *middleware.AuthMiddleware) {
	uploads := r.Group("/uploads", auth.RequireAuth())
	uploads.POST("/product-image", h.Upload.ProductImageUpload())

}
