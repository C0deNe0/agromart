package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/repository"
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
		&product.ProductCreateInput{},
		func(c echo.Context, req *product.ProductCreateInput) (*product.ProductResponse, error) {

			userID := middleware.GetUserID(c)

			if req.CompanyID == uuid.Nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "company id is required")
			}
			input := product.ProductCreateInput{
				CompanyID:  req.CompanyID,
				Name:       req.Name,
				CategoryID: req.CategoryID,
				Unit:       req.Unit,
				Origin:     req.Origin,
				Price:      req.Price,
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

func (h *ProductHandler) GetProductByID() echo.HandlerFunc {
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

func (h *ProductHandler) ListProducts() echo.HandlerFunc {
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

func (h *ProductHandler) UpdateProduct() echo.HandlerFunc {
	return Handle(
		&product.ProductUpdateInput{},
		func(c echo.Context, req *product.ProductUpdateInput) (*product.ProductResponse, error) {
			userID := middleware.GetUserID(c)

			input := product.ProductUpdateInput{
				ID:          req.ID,
				Name:        req.Name,
				Description: req.Description,
				CategoryID:  req.CategoryID,
				Unit:        req.Unit,
				Origin:      req.Origin,
				Price:       req.Price,
				IsActive:    req.IsActive,
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

func (h *ProductHandler) DeleteProduct() echo.HandlerFunc {
	return HandleNoContent(
		&product.DeleteProductRequest{},
		func(c echo.Context, req *product.DeleteProductRequest) error {
			userID := middleware.GetUserID(c)
			return h.productService.Delete(c.Request().Context(), userID, req.ID)
		},
		http.StatusNoContent,
	)
}
