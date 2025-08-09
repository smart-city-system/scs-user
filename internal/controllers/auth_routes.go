package http

import (
	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/login", h.Login())
}
