package handlers

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/middleware"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

//We do NOT use Handle[...] here on purpose
// This keeps /me dead simple and easy to debug.

func (h *UserHandler) Me() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := middleware.GetUserID(c)

		resp, err := h.userService.GetMe(c.Request().Context(), userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, resp)
	}
}
