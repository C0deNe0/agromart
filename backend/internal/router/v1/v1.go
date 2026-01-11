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
	authRoutes.POST("/logout", h.Auth.Logout())
	authRoutes.POST("/refresh", h.Auth.Refresh())
	//googleLogin
	authRoutes.POST("/google/login", h.Auth.LoginWithGoogleIDToken())

	//----PROTECTED ROUTES
	api := r.Group("")
	api.Use(auth.RequireAuth())

	//USER
	api.GET("/user/me", h.User.Me())

	//COMPANIES
	RegisterCompanyRoutes(api, h, auth)

	//products
	RegisterProductRoutes(api, h)

	//admin
	RegisterAdminRoutes(api, h, admin)

}
