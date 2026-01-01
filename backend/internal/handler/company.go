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
	service *service.CompanyService
}

func NewCompanyHandler(service *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		service: service,
	}
}

func (h *CompanyHandler) Create() echo.HandlerFunc {
	return Handle(
		&company.CreateCompanyRequest{},
		func(c echo.Context, req *company.CreateCompanyRequest) (*company.Company, error) {
			userID := middleware.GetUserID(c)
			input := company.CreateCompanyInput{
				Name:          req.Name,
				Description:   req.Description,
				LogoURL:       req.LogoURL,
				BusinessEmail: req.BusinessEmail,
				BusinessPhone: req.BusinessPhone,
				City:          req.City,
				State:         req.State,
				Pincode:       req.Pincode,
			}

			return h.service.Create(c.Request().Context(), userID, input)
		},
		http.StatusCreated,
	)
}

func (h *CompanyHandler) ListMine() echo.HandlerFunc {
	return func(c echo.Context) error {

		ownerID := middleware.GetUserID(c)
		result, err := h.service.ListByOwnerID(c.Request().Context(), ownerID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, result)
	}
}

func(h *CompanyHandler) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		companyID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		adminID := middleware.GetUserID(c)

		if err := h.service.DeleteCompany(c.Request().Context(), adminID, companyID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, "Company deleted successfully")
	}
}


func (h *CompanyHandler) Approve() echo.HandlerFunc {

	return func(c echo.Context) error {
		adminID := middleware.GetUserID(c)
		companyID, _ := uuid.Parse(c.Param("id"))

		return h.service.ApproveCompany(c.Request().Context(), adminID, companyID)
	}
}

func (h *CompanyHandler) Reject() echo.HandlerFunc {
	return func(c echo.Context) error {
		companyID, _ := uuid.Parse(c.Param("id"))

		return h.service.RejectCompany(c.Request().Context(), companyID)
	}
}

func (h *CompanyHandler) ListPending() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := h.service.ListPending(c.Request().Context())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, result)
	}
}
