package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CompanyHandler struct {
	companyService service.CompanyService
}

func NewCompanyHandler(companyService service.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

func (h *CompanyHandler) Create(c echo.Context) echo.HandlerFunc {
	return Handle(
		&company.CreateCompanyInput{},
		func(c echo.Context, req *company.CreateCompanyInput) (*company.Company, error) {
			userID := middleware.GetUserID(c)
			return h.companyService.Create(c.Request().Context(), userID, *req)
		},
		http.StatusCreated,
	)
}

func (h *CompanyHandler) ListByOwnerID(c echo.Context) error {
	ownerID := middleware.GetUserID(c)
	result, err := h.companyService.ListByOwnerID(c.Request().Context(), ownerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, result)
}

func (h *CompanyHandler) ListPending(c echo.Context) error {
	result, err := h.companyService.ListPending(c.Request().Context())
	if err != nil {

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}
	return c.JSON(http.StatusOK, result)
}

func (h *CompanyHandler) Approve(c echo.Context) error {
	type Request struct {
		ID uuid.UUID `param:"id" validate:"required"`
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	adminID := middleware.GetUserID(c)
	if err := h.companyService.ApproveCompany(c.Request().Context(), adminID, req.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Company approved successfully")
}

func (h *CompanyHandler) Reject(c echo.Context) error {
	type Request struct {
		ID uuid.UUID `param:"id" validate:"required"`
	}

	var req Request
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.companyService.RejectCompany(c.Request().Context(), req.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Company rejected!!")
}
