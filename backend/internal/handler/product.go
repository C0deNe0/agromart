package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/repository/productRepo"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/google/uuid"

	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/labstack/echo/v4"
)

//take care of hte middleware.GetUserID(c) => remporary fix with c.Get("user_id").(uuid.UUID)

type ProductHandler struct {
	Handler
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) CreateProduct() echo.HandlerFunc {
	return Handle(
		&product.CreateProductRequest{},
		func(c echo.Context, req *product.CreateProductRequest) (*product.ProductResponse, error) {

			userID := middleware.GetUserID(c)

			created, err := h.productService.Create(c.Request().Context(), userID, req)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return created, nil

		},
		http.StatusCreated,
	)

}

func (h *ProductHandler) GetProductByID() echo.HandlerFunc {
	return Handle(
		&product.GetProductByIDRequest{},
		func(c echo.Context, req *product.GetProductByIDRequest) (*product.ProductResponse, error) {
			userID := middleware.GetUserID(c)

			p, err := h.productService.GetByID(c.Request().Context(), req.ID, &userID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return p, nil
		},
		http.StatusOK,
	)
}

func (h *ProductHandler) ListProducts() echo.HandlerFunc {
	return Handle(
		&product.ListProductsQuery{},
		func(c echo.Context, req *product.ListProductsQuery) (interface{}, error) {
			var userID *uuid.UUID
			if id := middleware.GetUserID(c); id != uuid.Nil {
				userID = &id
			}

			filter := productRepo.ProductFilter{
				CompanyID:      req.CompanyID,
				CategoryID:     req.CategoryID,
				Search:         req.Search,
				ApprovalStatus: req.ApprovalStatus,
				IsActive:       req.IsActive,
				Page:           req.Page,
				Limit:          req.Limit,
			}
			result, err := h.productService.List(c.Request().Context(), userID, filter)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return result, nil
		},

		http.StatusOK,
	)
}

func (h *ProductHandler) UpdateProduct() echo.HandlerFunc {
	return Handle(
		&product.UpdateProductRequest{},
		func(c echo.Context, req *product.UpdateProductRequest) (*product.ProductResponse, error) {
			userID := middleware.GetUserID(c)

			updated, err := h.productService.Update(c.Request().Context(), userID, req.ID, req)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			p, err := h.productService.GetByID(c.Request().Context(), updated.ID, &userID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return p, nil
		},
		http.StatusOK,
	)
}

func (h *ProductHandler) DeleteProduct() echo.HandlerFunc {
	return HandleNoContent(
		&product.DeleteProductRequest{},
		func(c echo.Context, req *product.DeleteProductRequest) error {
			userID := middleware.GetUserID(c)
			err := h.productService.Delete(c.Request().Context(), userID, req.ID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			return nil
		},
		http.StatusNoContent,
	)
}

func (h *ProductHandler) ResubmitProduct() echo.HandlerFunc {
	return Handle(
		&product.ResubmitProductRequest{},
		func(c echo.Context, req *product.ResubmitProductRequest) (interface{}, error) {
			userID := middleware.GetUserID(c)

			err := h.productService.Resubmit(c.Request().Context(), userID, req.ProductID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return map[string]string{
				"message": "Product resubmitted for approval successfully",
			}, nil
		},
		http.StatusOK,
	)
}

func (h *ProductHandler) GetApprovalHistory() echo.HandlerFunc {
	return Handle(
		&product.GetProductByIDRequest{},
		func(c echo.Context, req *product.GetProductByIDRequest) (interface{}, error) {
			history, err := h.productService.GetApprovalHistory(c.Request().Context(), req.ID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return history, nil
		},
		http.StatusOK,
	)
}

//IMAGE HANDLERS

func (h *ProductHandler) GenerateImageUploadURL() echo.HandlerFunc {
	return Handle(
		&product.GenerateImageUploadURLRequest{},
		func(c echo.Context, req *product.GenerateImageUploadURLRequest) (*product.ImageUploadURLResponse, error) {
			userID := middleware.GetUserID(c)

			response, err := h.productService.GenerateImageUploadURL(c.Request().Context(), userID, req)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return response, nil
		},
		http.StatusOK,
	)
}

func (h *ProductHandler) DeleteImage() echo.HandlerFunc {
	return HandleNoContent(
		&product.DeleteProductImageRequest{},
		func(c echo.Context, req *product.DeleteProductImageRequest) error {
			userID := middleware.GetUserID(c)

			err := h.productService.DeleteImage(c.Request().Context(), userID, req.ProductID, req.ImageID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return nil
		},
		http.StatusNoContent,
	)
}

func (h *ProductHandler) SetPrimaryImage() echo.HandlerFunc {
	return Handle(
		&product.DeleteProductImageRequest{}, // Reuse same structure
		func(c echo.Context, req *product.DeleteProductImageRequest) (interface{}, error) {
			userID := middleware.GetUserID(c)

			err := h.productService.SetPrimaryImage(c.Request().Context(), userID, req.ProductID, req.ImageID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return map[string]string{
				"message": "Primary image set successfully",
			}, nil
		},
		http.StatusOK,
	)
}

// =============================================
// VARIANT MANAGEMENT
// =============================================

func (h *ProductHandler) CreateVariant() echo.HandlerFunc {
	return Handle(
		&product.CreateVariantRequest{},
		func(c echo.Context, req *product.CreateVariantRequest) (*product.ProductVariantResponse, error) {
			userID := middleware.GetUserID(c)

			variant, err := h.productService.CreateVariant(c.Request().Context(), userID, req)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return &product.ProductVariantResponse{
				ID:                variant.ID,
				ProductID:         variant.ProductID,
				Label:             variant.Label,
				QuantityValue:     variant.QuantityValue,
				QuantityUnit:      variant.QuantityUnit,
				Price:             variant.Price,
				StockQuantity:     variant.StockQuantity,
				LowStockThreshold: variant.LowStockThreshold,
				IsLowStock:        variant.IsLowStock(),
				IsAvailable:       variant.IsAvailable,
				CreatedAt:         variant.CreatedAt,
				UpdatedAt:         variant.UpdatedAt,
			}, nil
		},
		http.StatusCreated,
	)
}

func (h *ProductHandler) UpdateVariant() echo.HandlerFunc {
	return Handle(
		&product.UpdateVariantRequest{},
		func(c echo.Context, req *product.UpdateVariantRequest) (*product.ProductVariantResponse, error) {
			userID := middleware.GetUserID(c)

			variant, err := h.productService.UpdateVariant(c.Request().Context(), userID, req.ProductID, req.VariantID, req)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return &product.ProductVariantResponse{
				ID:        variant.ID,
				ProductID: variant.ProductID,

				Label:             variant.Label,
				QuantityValue:     variant.QuantityValue,
				QuantityUnit:      variant.QuantityUnit,
				Price:             variant.Price,
				StockQuantity:     variant.StockQuantity,
				LowStockThreshold: variant.LowStockThreshold,
				IsLowStock:        variant.IsLowStock(),
				IsAvailable:       variant.IsAvailable,
				CreatedAt:         variant.CreatedAt,
				UpdatedAt:         variant.UpdatedAt,
			}, nil
		},
		http.StatusOK,
	)
}

func (h *ProductHandler) DeleteVariant() echo.HandlerFunc {
	return HandleNoContent(
		&product.DeleteVariantRequest{},
		func(c echo.Context, req *product.DeleteVariantRequest) error {
			userID := middleware.GetUserID(c)

			err := h.productService.DeleteVariant(c.Request().Context(), userID, req.ProductID, req.VariantID)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return nil
		},
		http.StatusNoContent,
	)
}
