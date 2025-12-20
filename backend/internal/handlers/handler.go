package handlers

import (
	"github.com/C0deNe0/agromart/internal/validation"
	"github.com/labstack/echo/v4"
)

type Handler struct{}
type HandlerFunc[Req validation.Validatable, Res any] func(
	c echo.Context,
	req Req,
) (Res, error)

func handleRequest[Req validation.Validatable](
	c echo.Context,
	req Req,
	handler func(c echo.Context, req Req) (any, error),
	status int,
) error {

	//binding the request and validating
	if err := validation.BindAndValidate(c, req); err != nil {
		return c.JSON(status, echo.Map{
			"error": err.Error(),
		})
	}

	//executing the handler
	result, err := handler(c, req)
	if err != nil {
		return c.JSON(status, echo.Map{
			"error": err.Error(),
		})
	}

	if result == nil {
		return c.NoContent(status)
	}

	return c.JSON(status, result)
}

func Handle[Req validation.Validatable, Res any](
	req Req,
	handler HandlerFunc[Req, Res],
	status int,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (any, error) {
			return handler(c, req)
		},
			status,
		)
	}
}

func HandleNoContent[Req validation.Validatable](
	req Req,
	handler func(c echo.Context, req Req) error,
	status int,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		return handleRequest(c, req, func(c echo.Context, req Req) (any, error) {
			return nil, handler(c, req)
		},
			status,
		)
	}
}

// 1. Read request data
// 2. Validate it
// 3. Call your logic
// 4. Send response
