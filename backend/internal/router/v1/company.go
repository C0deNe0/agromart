package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterCompanyRoutes(r *echo.Group, h *handler.Handlers, auth *middleware.AuthMiddleware) {
	company := r.Group("/companies")

	company.GET("", h.Company.ListCompanies())      //all companies
	company.GET("/:id", h.Company.GetCompanyByID()) //by id

	company.POST("", h.Company.CreateCompany())
	company.PUT("/:id", h.Company.UpdateCompany())
	company.DELETE("/:id", h.Company.DeleteCompany())
	company.POST("/:id/resubmit", h.Company.ResubmitCompany())

	company.GET("/:id/histroy", h.Company.GetApprovalHistory()) // all approval histroy

	company.POST("/:id/follow", h.Company.FollowCompany())
	company.DELETE("/:id/unfollow", h.Company.UnfollowCompany())
	company.GET("/:id/follow-status", h.Company.GetFollowStatus())
	company.GET("/:id/followers", h.Company.ListFollowers())

	company.GET("/followed/me", h.Company.ListFollowedCompanies())
}
