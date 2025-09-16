package server

import (
	"net/http"
	controller "scs-user/internal/controllers"
	my_middleware "scs-user/internal/middlewares"
	repository "scs-user/internal/repositories"
	service "scs-user/internal/services"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func (s *Server) MapHandlers(e *echo.Echo) error {
	// Init repositories
	userRepo := repository.NewUserRepository(s.db)
	userPremiseRepo := repository.NewUserPremiseRepository(s.db)

	// Init service
	userService := service.NewUserService(*userRepo, *userPremiseRepo, *s.producer)
	authService := service.NewAuthService(*userRepo)
	// Init handlers
	userHandler := controller.NewUserHandler(*userService)
	authHandler := controller.NewAuthHandler(*authService)

	// Enable CORS for all origins
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: false,
	}))

	mw := my_middleware.NewMiddlewareManager(s.cfg, []string{"*"}, s.logger)
	e.Use(mw.RequestLoggerMiddleware)
	e.Use(mw.ErrorHandlerMiddleware)
	e.Use(mw.ResponseStandardizer)
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
