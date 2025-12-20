package handlers

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/google/uuid"

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

func (h *ProductHandler) CreateProduct(c echo.Context) echo.HandlerFunc {
	return Handle(
		&product.CreateProductRequest{},
		func(c echo.Context, req *product.CreateProductRequest) (*product.ProductResponse, error) {

			userID := c.Get("user_id").(uuid.UUID)

			if req.CompanyID == nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "company id is required")
			}
			input := service.ProductCreateInput{
				CompanyID:    *req.CompanyID,
				Name:         req.Name,
				CategoryID:   req.CategoryID,
				Unit:         req.Unit,
				Origin:       req.Origin,
				PriceDisplay: req.PriceDisplay,
			}
			p, err := h.productService.Create(
				c.Request().Context(),
				userID,
				input,
			)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return product.ToProductResponse(p), nil
		},
		http.StatusCreated,
	)

}

func (h *ProductHandler) GetProductByID(c echo.Context) echo.HandlerFunc {
	return Handle(
		&product.GetProductByIDRequest{},
		func(c echo.Context, req *product.GetProductByIDRequest) (*product.ProductResponse, error) {
			p, err := h.productService.GetByID(c.Request().Context(), req.ID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return product.ToProductResponse(p), nil
		},
		http.StatusOK,
	)
}

func (h *ProductHandler) ListProducts(c echo.Context) echo.HandlerFunc {
	return Handle(
		&product.ListProductsQuery{},
		func(c echo.Context, req *product.ListProductsQuery) (*model.PaginatedResponse[product.ProductResponse], error) {
			filter := repository.ProductFilter{
				CompanyID:  req.CompanyID,
				CategoryID: req.CategoryID,
				Search:     req.Search,
				IsActive:   req.IsActive,
				Page:       *req.Page,
				Limit:      *req.Limit,
			}
			result, err := h.productService.List(c.Request().Context(), filter)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return product.MapProductPage(result), nil
		},

		http.StatusOK,
	)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) echo.HandlerFunc {
	return Handle(
		&product.UpdateProductRequest{},
		func(c echo.Context, req *product.UpdateProductRequest) (*product.ProductResponse, error) {
			userID := c.Get("user_id").(uuid.UUID)

			input := service.ProductUpdateInput{
				ID:           req.ID,
				Name:         req.Name,
				Description:  req.Description,
				CategoryID:   req.CategoryID,
				Unit:         req.Unit,
				Origin:       req.Origin,
				PriceDisplay: req.PriceDisplay,
				IsActive:     req.IsActive,
			}

			p, err := h.productService.Update(c.Request().Context(), userID, input)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return product.ToProductResponse(p), nil
		},
		http.StatusOK,
	)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) echo.HandlerFunc {
	return HandleNoContent(
		&product.DeleteProductRequest{},
		func(c echo.Context, req *product.DeleteProductRequest) error {
			userID := c.Get("user_id").(uuid.UUID)
			return h.productService.Delete(c.Request().Context(), userID, req.ID)
		},
		http.StatusNoContent,
	)
}
