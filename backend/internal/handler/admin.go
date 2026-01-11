package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	Handler
	companyService *service.CompanyService
	productService *service.ProductService
}

func NewAdminHandler(companyService *service.CompanyService, productService *service.ProductService) *AdminHandler {
	return &AdminHandler{
		companyService: companyService,
		productService: productService,
	}
}

func (h *AdminHandler) ApproveCompany() echo.HandlerFunc {
	return Handle(
		&company.ApproveCompanyRequest{},
		func(c echo.Context, req *company.ApproveCompanyRequest) (interface{}, error) {
			adminID := middleware.GetUserID(c)

			err := h.companyService.Approve(c.Request().Context(), req.CompanyID, adminID, req.Notes)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return map[string]string{
				"message": "Company approved successfully",
			}, nil
		},
		http.StatusOK,
	)
}

func (h *AdminHandler) RejectCompany() echo.HandlerFunc {
	return Handle(
		&company.RejectCompanyRequest{},
		func(c echo.Context, req *company.RejectCompanyRequest) (interface{}, error) {
			adminID := middleware.GetUserID(c)

			err := h.companyService.Reject(c.Request().Context(), req.CompanyID, adminID, req.Reason, req.Notes)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return map[string]string{
				"message": "Company rejected successfully",
			}, nil
		},
		http.StatusOK,
	)
}

func (h *AdminHandler) CountPendingCompanyApprovals() echo.HandlerFunc {
	return Handle(
		&company.CountPendingApprovalsRequest{},
		func(c echo.Context, req *company.CountPendingApprovalsRequest) (interface{}, error) {
			count, err := h.companyService.CountPendingApprovals(c.Request().Context())
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return map[string]int{
				"pendingCount": count,
			}, nil
		},
		http.StatusOK,
	)
}

// =============================================
// PRODUCT APPROVAL
// =============================================

func (h *AdminHandler) ApproveProduct() echo.HandlerFunc {
	return Handle(
		&product.ApproveProductRequest{},
		func(c echo.Context, req *product.ApproveProductRequest) (interface{}, error) {
			adminID := middleware.GetUserID(c)

			err := h.productService.Approve(c.Request().Context(), req.ProductID, adminID, req.Notes)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return map[string]string{
				"message": "Product approved successfully",
			}, nil
		},
		http.StatusOK,
	)
}

func (h *AdminHandler) RejectProduct() echo.HandlerFunc {
	return Handle(
		&product.RejectProductRequest{},
		func(c echo.Context, req *product.RejectProductRequest) (interface{}, error) {
			adminID := middleware.GetUserID(c)

			err := h.productService.Reject(c.Request().Context(), req.ProductID, adminID, req.Reason, req.Notes)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return map[string]string{
				"message": "Product rejected successfully",
			}, nil
		},
		http.StatusOK,
	)
}

func (h *AdminHandler) CountPendingProducts() echo.HandlerFunc {
	return Handle(
		&product.CountPendingApprovalsRequest{},
		func(c echo.Context, req *product.CountPendingApprovalsRequest) (interface{}, error) {
			count, err := h.productService.CountPendingApprovals(c.Request().Context())
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return map[string]int{
				"pendingCount": count,
			}, nil
		},
		http.StatusOK,
	)
}
