package v1

import (
	"github.com/C0deNe0/agromart/internal/handler"
	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/labstack/echo/v4"
)

func RegisterV1Routes(c echo.Context, r *echo.Group, h *handler.Handlers, auth *middleware.AuthMiddleware, admin *middleware.AdminMiddleware) {

	//----REGISTERING THE USER ROUTES
	authRoutes := r.Group("/auth")
	authRoutes.POST("/register", h.Auth.Register())
	authRoutes.POST("/login", h.Auth.Login())
	authRoutes.POST("/refresh", h.Auth.Refresh())

	//google auth
	authRoutes.GET("/google/login", h.Auth.GoogleLogin())
	authRoutes.GET("/google/callback", h.Auth.GoogleCallback())

	//----PROTECTED ROUTES
	api := r.Group("")
	api.Use(auth.RequireAuth())

	//USER
	api.GET("/user/me", h.User.Me())

	//COMPANIES
	RegisterCompanyRoutes(c, api, h)

	//products
	RegisterProductRoutes(c, api, h)

	//admin
	RegisterAdminRoutes(api, h, admin)

}
