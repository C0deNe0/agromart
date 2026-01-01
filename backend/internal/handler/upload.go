package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/C0deNe0/agromart/internal/model/upload"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	UploadService *service.UploadService
}

func NewUploadHandler(uploadService *service.UploadService) *UploadHandler {
	return &UploadHandler{
		UploadService: uploadService,
	}
}

func (h *UploadHandler) ProductImageUpload() echo.HandlerFunc {
	return func(c echo.Context) error {

		userID := middleware.GetUserID(c)

		//ye check kar
		productID, err := uuid.Parse(c.QueryParam("product_id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid product id")
		}

		contentType := c.QueryParam("content_type")
		if contentType == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "content type is required")
		}

		req := upload.UploadRequest{
			Type:        upload.UploadProductImage,
			UserID:      userID,
			ProductID:   &productID,
			ContentType: contentType,
		}

		resp, err := h.UploadService.PresignUpload(c.Request().Context(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusForbidden, err.Error())
		}

		return c.JSON(http.StatusOK, resp)

	}
}
