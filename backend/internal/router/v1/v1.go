package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(c echo.Context, r *echo.Group, h *handler.Handlers, tokenManager *utils.TokenManager) {
	//----REGISTERING THE USER ROUTES
	auth := r.Group("/auth")
	auth.POST("/register", h.Auth.Register())
	auth.POST("/login", h.Auth.Login())

	//google auth
	auth.GET("/google/login", h.Auth.GoogleLogin())
	auth.GET("/google/callback", h.Auth.GoogleCallback())

	//----PROTECTED ROUTES
	api := r.Group("")
	api.Use(middleware.AuthMiddleware(tokenManager))

	//USER
	api.GET("/user", h.User.Me())

	//COMPANIES
	RegisterCompanyRoutes(c, api, h)

	//products
	RegisterProductRoutes(c, api, h)

	//admin
	admin := api.Group("/admin")
	admin.Use(middleware.AdminOnly)

	RegisterAdminRoutes(admin, h)

}
