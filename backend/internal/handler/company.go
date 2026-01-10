package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CompanyHandler struct {
	Handler
	companyService *service.CompanyService
}

func NewCompanyHandler(companyService *service.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

func (h *CompanyHandler) CreateCompany() echo.HandlerFunc {
	return Handle(
		&company.CreateCompanyRequest{},
		func(c echo.Context, req *company.CreateCompanyRequest) (*company.CompanyResponse, error) {
			userID := middleware.GetUserID(c)
			comp := company.Company{
				Name:          req.Name,
				Description:   req.Description,
				LogoURL:       req.LogoURL,
				BusinessEmail: req.BusinessEmail,
				BusinessPhone: req.BusinessPhone,

				City:      req.City,
				State:     req.State,
				Pincode:   req.Pincode,
				GSTNumber: req.GSTNumber,
				PANNumber: req.PANNumber,
			}

			if req.ProductVisibility != nil {
				comp.ProductVisibility = *req.ProductVisibility
			}

			created, err := h.companyService.Create(c.Request().Context(), userID, comp)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}
			return company.ToCompanyResponse(created, nil), nil
		},

		http.StatusCreated,
	)
}

func (h *CompanyHandler) GetCompanyByID() echo.HandlerFunc {
	return Handle(
		&company.GetCompanyByIDRequest{},
		func(c echo.Context, req *company.GetCompanyByIDRequest) (*company.CompanyResponse, error) {
			var userID *uuid.UUID
			if id := middleware.GetUserID(c); id != uuid.Nil {
				userID = &id

			}

			comp, err := h.companyService.GetByID(c.Request().Context(), req.ID, userID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return comp, nil
		},
		http.StatusOK,
	)
}

func (h *CompanyHandler) ListCompanies() echo.HandlerFunc {
	return Handle(
		&company.ListCompanyQuery{},
		func(c echo.Context, req *company.ListCompanyQuery) (interface{}, error) {
			var userID *uuid.UUID
			if id := middleware.GetUserID(c); id != uuid.Nil {
				userID = &id
			}
			filter := repository.CompanyFilter{
				OwnerID:        req.OwnerID,
				Search:         req.Search,
				ApprovalStatus: req.ApprovalStatus,
				IsActive:       req.IsActive,
				Page:           req.Page,
				Limit:          req.Limit,
			}
			return h.companyService.List(c.Request().Context(), userID, filter)
		},
		http.StatusOK,
	)
}

func (h *CompanyHandler) UpdateCompany() echo.HandlerFunc {
	return Handle(
		&company.UpdateCompanyRequest{},
		func(c echo.Context, req *company.UpdateCompanyRequest) (*company.CompanyResponse, error) {
			userID := middleware.GetUserID(c)
			updated, err := h.companyService.Update(c.Request().Context(), userID, req.ID, req)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return company.ToCompanyResponse(updated, nil), nil
		},
		http.StatusOK,
	)
}

func (h *CompanyHandler) DeleteCompany() echo.HandlerFunc {
	return HandleNoContent(
		&company.DeleteCompanyRequest{},
		func(c echo.Context, req *company.DeleteCompanyRequest) error {
			userID := middleware.GetUserID(c)

			err := h.companyService.Delete(c.Request().Context(), userID, req.ID)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}
			return nil
		},
		http.StatusNoContent,
	)
}

func (h *CompanyHandler) ResubmitCompany() echo.HandlerFunc {
	return Handle(
		&company.ResubmitCompanyRequest{},
		func(c echo.Context, req *company.ResubmitCompanyRequest) (interface{}, error) {
			userID := middleware.GetUserID(c)

			err := h.companyService.Resubmit(c.Request().Context(), userID, req.ID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return map[string]string{
				"message": "Company resubmitted for approval successfully",
			}, nil
		},
		http.StatusOK,
	)
}

func (h *CompanyHandler) GetApprovalHistory() echo.HandlerFunc {
	return Handle(
		&company.GetApprovalHistoryRequest{},
		func(c echo.Context, req *company.GetApprovalHistoryRequest) (interface{}, error) {
			history, err := h.companyService.GetApprovalHistory(c.Request().Context(), req.ID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return history, nil
		},
		http.StatusOK,
	)
}

// =============================================
// FOLLOW COMPANY
// =============================================

func (h *CompanyHandler) FollowCompany() echo.HandlerFunc {
	return Handle(
		&company.FollowCompanyRequest{},
		func(c echo.Context, req *company.FollowCompanyRequest) (interface{}, error) {
			userID := middleware.GetUserID(c)

			err := h.companyService.Follow(c.Request().Context(), req.CompanyID, userID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return map[string]string{
				"message": "Successfully followed company",
			}, nil
		},
		http.StatusOK,
	)
}

// =============================================
// UNFOLLOW COMPANY
// =============================================

func (h *CompanyHandler) UnfollowCompany() echo.HandlerFunc {
	return Handle(
		&company.UnfollowCompanyRequest{},
		func(c echo.Context, req *company.UnfollowCompanyRequest) (interface{}, error) {
			userID := middleware.GetUserID(c)

			err := h.companyService.Unfollow(c.Request().Context(), req.CompanyID, userID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			return map[string]string{
				"message": "Successfully unfollowed company",
			}, nil
		},
		http.StatusOK,
	)
}

// =============================================
// GET FOLLOW STATUS
// =============================================

func (h *CompanyHandler) GetFollowStatus() echo.HandlerFunc {
	return Handle(
		&company.FollowCompanyRequest{},
		func(c echo.Context, req *company.FollowCompanyRequest) (*company.FollowStatusResponse, error) {
			userID := middleware.GetUserID(c)

			status, err := h.companyService.GetFollowStatus(c.Request().Context(), req.CompanyID, userID)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return status, nil
		},
		http.StatusOK,
	)
}

// =============================================
// LIST FOLLOWERS
// =============================================

func (h *CompanyHandler) ListFollowers() echo.HandlerFunc {
	return Handle(
		&company.ListFollowersQuery{},
		func(c echo.Context, req *company.ListFollowersQuery) (interface{}, error) {
			followers, err := h.companyService.ListFollowers(
				c.Request().Context(),
				req.CompanyID,
				req.Page,
				req.Limit,
			)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return followers, nil
		},
		http.StatusOK,
	)
}

// =============================================
// LIST FOLLOWED COMPANIES (MY FOLLOWED)
// =============================================

func (h *CompanyHandler) ListFollowedCompanies() echo.HandlerFunc {
	return Handle(
		&company.ListFollowedCompaniesQuery{},
		func(c echo.Context, req *company.ListFollowedCompaniesQuery) (interface{}, error) {
			userID := middleware.GetUserID(c)

			companies, err := h.companyService.ListFollowedCompanies(
				c.Request().Context(),
				userID,
				req.Page,
				req.Limit,
			)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return companies, nil
		},
		http.StatusOK,
	)
}
