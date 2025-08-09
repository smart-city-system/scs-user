package server

import (
	"context"
	"net/http"
	config "scs-user/config"
	logger "scs-user/pkg/logger"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

type Server struct {
	Echo   *echo.Echo
	cfg    *config.Config
	db     *gorm.DB
	logger logger.Logger
}

func NewServer(cfg *config.Config, db *gorm.DB, logger logger.Logger) *Server {
	return &Server{cfg: cfg, db: db, logger: logger, Echo: echo.New()}
}
func (s *Server) Run() error {
	// Map handlers
	if err := s.MapHandlers(s.Echo); err != nil {
		return err
	}

	// create http server
	server := &http.Server{
		Addr:           ":" + s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
	return s.Echo.StartServer(server)
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}
