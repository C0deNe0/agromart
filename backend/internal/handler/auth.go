package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/model/auth"
	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register() echo.HandlerFunc {
	return Handle(
		&user.RegisterRequest{},
		func(c echo.Context, req *user.RegisterRequest) (*user.AuthResponse, error) {
			return h.authService.RegisterWithEmail(
				c.Request().Context(),
				req.Email,
				req.Password,
				req.Name,
			)
		},
		http.StatusCreated,
	)
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return Handle(
		&user.LoginRequest{},
		func(c echo.Context, req *user.LoginRequest) (*user.AuthResponse, error) {
			return h.authService.LoginWithEmail(
				c.Request().Context(),
				req.Email,
				req.Password,
			)
		},
		http.StatusOK,
	)

}

func (h *AuthHandler) Refresh() echo.HandlerFunc {
	return Handle(
		&auth.RefreshRequest{},
		func(c echo.Context, req *auth.RefreshRequest) (*user.AuthResponse, error) {
			resp, err := h.authService.Refresh(
				c.Request().Context(),
				req.RefreshToken,
			)
			if err != nil {
				return nil, echo.NewHTTPError(
					http.StatusUnauthorized,
					"invalid or expired refresh token",
				)
			}
			return resp, nil
		},
		http.StatusOK,
	)
}

func (h *AuthHandler) Logout() echo.HandlerFunc {
	return Handle(
		&auth.LogoutRequest{},
		func(c echo.Context, req *auth.LogoutRequest) (map[string]interface{}, error) {
			if err := h.authService.Logout(
				c.Request().Context(),
				req.RefreshToken,
			); err != nil {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired refresh token")
			}
			return map[string]interface{}{
				"message": "logout successful",
			}, nil
		},
		http.StatusOK,
	)
}
func (h *AuthHandler) LoginWithGoogleIDToken() echo.HandlerFunc {
	return Handle(
		&auth.GoogleIDTokenRequest{},
		func(c echo.Context, req *auth.GoogleIDTokenRequest) (*user.AuthResponse, error) {

			googleUserClaims, err := utils.VerifyGoogleIDToken(
				c.Request().Context(),
				req.IDToken,
			)
			if err != nil || googleUserClaims.Email == "" {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid Google ID Token"+err.Error())
			}

			// 2. Call the service layer to handle login/registration using the verified claims
			resp, err := h.authService.LoginWithGoogle(
				c.Request().Context(),
				googleUserClaims.Sub,
				googleUserClaims.Email,
				googleUserClaims.Name,
				&googleUserClaims.Picture,
			)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, "authentication failed")
			}

			return resp, nil
		},
		http.StatusOK,
	)
}

//REGISTER----------

// HTTP → Handler → AuthService
//      → UserRepo.Create
//      → UserAuthMethodRepo.Create
//      → TokenManager.Generate
//      →

// LOGIN------------

// HTTP → Handler → AuthService
//      → UserAuthMethodRepo.GetByEmail
//      → VerifyPassword
//      → TokenManager.Generate
//      → Response

// GOOGLE AUTH
// HTTP → Handler → AuthService
//      → GoogleOAuth.Exchange
//      → FetchGoogleUser
//      → AuthService.LoginWithGoogle
//      → TokenManager.Generate
//      → Response

// func (h *AuthHandler) GoogleLogin() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		state := uuid.New().String()

// 		//store state in HTTP only cookie for 5 min

// 		c.SetCookie(&http.Cookie{
// 			Name:     "oauth_state",
// 			Value:    state,
// 			Path:     "/",
// 			HttpOnly: true,
// 			Secure:   true,
// 			MaxAge:   300,
// 			// Expires:  time.Now().Add(5 * time.Minute),
// 			SameSite: http.SameSiteLaxMode,
// 		})

// 		return c.Redirect(
// 			http.StatusTemporaryRedirect,
// 			h.googleOAuth.AuthURL(state),
// 		)
// 	}
// }

// func (h *AuthHandler) GoogleCallback() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		state := c.QueryParam("state")
// 		if state == "" {
// 			return echo.NewHTTPError(http.StatusBadRequest, "state is required")
// 		}
// 		code := c.QueryParam("code")
// 		if code == "" {
// 			return echo.NewHTTPError(http.StatusBadRequest, "code is required")
// 		}

// 		//validate state
// 		cookie, err := c.Cookie("oauth_state")
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusBadRequest, "state is required")
// 		}
// 		if cookie.Value != state {
// 			return echo.NewHTTPError(http.StatusBadRequest, "state does not match")
// 		}

// 		//clearing after use
// 		c.SetCookie(&http.Cookie{
// 			Name:     "oauth_state",
// 			Value:    "",
// 			Path:     "/",
// 			HttpOnly: true,
// 			Secure:   true,
// 			MaxAge:   -1,
// 			SameSite: http.SameSiteLaxMode,
// 		})

// 		token, err := h.googleOAuth.Exchange(c.Request().Context(), code)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 		}

// 		googleUser, err := utils.FetchGoogleUser(token.AccessToken)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 		}

// 		resp, err := h.authService.LoginWithGoogle(c.Request().Context(),
// 			googleUser.Sub,
// 			googleUser.Email,
// 			googleUser.Name,
// 			&googleUser.Picture,
// 		)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 		}

// 		return c.JSON(http.StatusOK, resp)

// 	}
// }
