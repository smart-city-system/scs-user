package http

import (
	"scs-user/internal/dto"
	services "scs-user/internal/services"
	"scs-user/pkg/errors"
	"scs-user/pkg/validation"

	"github.com/labstack/echo/v4"
)

// Handler
type AuthHandler struct {
	svc services.AuthService
}

// NewHandler constructor
func NewAuthHandler(svc services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse "Login successful"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginReq dto.LoginRequest

		if err := c.Bind(&loginReq); err != nil {
			return err
		}
		// Validate the DTO
		if err := validation.ValidateStruct(loginReq); err != nil {
			return err
		}
		token, err := h.svc.Login(c.Request().Context(), &loginReq)
		if err != nil {
			return err
		}

		return c.JSON(200, token)
	}
}

// ValidateToken godoc
// @Summary Validate JWT token
// @Description Validate if a JWT token is valid and not expired
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ValidateTokenRequest true "Token to validate"
// @Success 200 {object} dto.ValidateTokenResponse "Token validation result"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Invalid token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/validate-token [post]
func (h *AuthHandler) ValidateToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		validateTokenDto := &dto.ValidateTokenRequest{}
		if err := c.Bind(validateTokenDto); err != nil {
			return errors.NewBadRequestError("Invalid request body")
		}
		// Validate the DTO
		if err := validation.ValidateStruct(validateTokenDto); err != nil {
			return err
		}

		result, err := h.svc.ValidateToken(c.Request().Context(), validateTokenDto.Token)
		if err != nil {
			return err
		}
		return c.JSON(200, result)
	}
}
