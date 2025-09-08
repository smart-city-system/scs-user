package http

import (
	middleware "scs-user/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func (h *UserHandler) RegisterRoutes(g *echo.Group, mw *middleware.MiddlewareManager) {

	g.POST("", mw.JWTAuth(h.CreateUser()))
	g.GET("", mw.JWTAuth(h.GetUsers()))
	g.GET("/me", mw.JWTAuth(h.GetMe()))

}
