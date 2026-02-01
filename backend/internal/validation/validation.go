package validation

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

var validate = validator.New()

type Validatable interface {
	Validate() error
}

func BindAndValidate(c echo.Context, req Validatable) error {
	if err := c.Bind(req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

