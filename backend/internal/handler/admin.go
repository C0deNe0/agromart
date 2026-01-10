package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	Handler
	companyService *service.CompanyService
}

func NewAdminHandler(companyService *service.CompanyService) *AdminHandler {
	return &AdminHandler{
		companyService: companyService,
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

func (h *AdminHandler) CountPendingApprovals() echo.HandlerFunc {
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
