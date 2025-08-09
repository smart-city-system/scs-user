package server

import (
	"net/http"
	controller "scs-user/internal/controllers"
	repository "scs-user/internal/repositories"
	service "scs-user/internal/services"

	middleware "scs-user/internal/middlewares"

	"github.com/labstack/echo/v4"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init repositories
	userRepo := repository.NewUserRepository(s.db)

	// Init service
	userService := service.NewUserService(*userRepo)
	authService := service.NewUserService(*userRepo)
	// Init handlers
	userHandler := controller.NewUserHandler(*userService)
	authHandler := controller.NewAuthHandler(*authService)

	mw := middleware.NewMiddlewareManager(s.cfg, []string{"*"}, s.logger)
	e.Use(mw.RequestLoggerMiddleware)
	e.Use(mw.ErrorHandlerMiddleware)
	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	usersGroup := v1.Group("/users")
	authGroup := v1.Group("/auth")

	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})
	userHandler.RegisterRoutes(usersGroup, mw)
	authHandler.RegisterRoutes(authGroup)

	return nil

}
