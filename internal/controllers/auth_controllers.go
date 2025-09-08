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
