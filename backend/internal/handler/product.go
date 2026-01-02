package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/C0deNe0/agromart/internal/service"

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

			p := &product.Product{
				CompanyID:   req.CompanyID,
				CategoryID:  req.CategoryID,
				Name:        req.Name,
				Description: req.Description,
				Unit:        req.Unit,
				Origin:      req.Origin,
				Price:       req.Price,
				IsActive:    true,
			}

			created, err := h.productService.Create(c.Request().Context(), userID, p)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}


			return product.ToProductResponse(created), nil

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
				Page:       req.Page,
				Limit:      req.Limit,
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

func (h *ProductHandler) ListProductsWithCategory() echo.HandlerFunc {
    return Handle(
        &product.ListProductsQuery{},
        func(c echo.Context, req *product.ListProductsQuery) (*model.PaginatedResponse[product.ProductWithCategoryResponse], error) {
            filter := repository.ProductFilter{
                CompanyID:  req.CompanyID,
                CategoryID: req.CategoryID,
                Search:     req.Search,
                IsActive:   req.IsActive,
                Page:       req.Page,
                Limit:      req.Limit,
            }
            
            result, err := h.productService.ListWithCategory(c.Request().Context(), filter)
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
			return product.ToProductResponse(updated), nil
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
