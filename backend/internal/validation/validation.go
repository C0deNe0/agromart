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

	return ValidateStruct(req)
}

func ValidateStruct(v Validatable) error {
	if err := validate.Struct(v); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

// var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// func IsValidUUID(uuid string) bool {
// 	return uuidRegex.MatchString(uuid)
// }
