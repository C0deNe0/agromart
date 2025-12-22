package router

import (
	"github.com/C0deNe0/agromart/internal/handlers"
	"github.com/C0deNe0/agromart/internal/lib/utils"
	v1 "github.com/C0deNe0/agromart/internal/router/v1"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func NewRouter(h *handlers.Handlers, tokenManager *utils.TokenManager) *echo.Echo {
	e := echo.New()

	//--GLOBAL MIDDLEWARES
	e.Use(
		echoMiddleware.Recover(),
		echoMiddleware.RequestID(),
		echoMiddleware.CORS(),
		echoMiddleware.BodyLimit("10MB"),
	)

	//----REGISTERING THE SYSTEM ROUTES
	RegisterSystemRoutes(e, h)

	//----REGISTERING THE V1 ROUTES
	v1Route := e.Group("/v1")
	v1.RegisterV1Routes(v1Route, h, tokenManager)

	return e
}
