package http

import (
	"scs-user/internal/dto"
	services "scs-user/internal/services"
	"scs-user/pkg/validation"

	"github.com/labstack/echo/v4"
)

// Handler
type AuthHandler struct {
	svc services.UserService
}

// NewHandler constructor
func NewAuthHandler(svc services.UserService) *AuthHandler {
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
