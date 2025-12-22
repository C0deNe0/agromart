package handlers

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CompanyHandler struct {
	companyService service.CompanyService
}

func NewCompanyHandler(companyService service.CompanyService) CompanyHandler {
	return CompanyHandler{
		companyService: companyService,
	}
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
	if err := h.companyService.Approve(c.Request().Context(), adminID, req.ID); err != nil {
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

	adminID := middleware.GetUserID(c)
	if err := h.companyService.Reject(c.Request().Context(), adminID, req.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Company rejected!!")
}