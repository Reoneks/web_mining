package server

import (
	"context"
	"errors"
	"net/http"

	"dyploma/config"
	"dyploma/server/handlers"
	"dyploma/server/middleware"

	"github.com/labstack/echo/v4"
	log "github.com/rs/zerolog/log"
)

type HTTPServer struct {
	router *echo.Echo

	cfg      *config.Config
	handlers *handlers.Handler
}

func NewHTTPServer(
	cfg *config.Config,
	handlers *handlers.Handler,
) *HTTPServer {
	return &HTTPServer{
		cfg:      cfg,
		handlers: handlers,
	}
}

func (s *HTTPServer) Start(_ context.Context) error {
	s.router = echo.New()
	s.router.Use(middleware.LoggerMiddleware(), middleware.CorsMiddleware(), middleware.RecoverMiddleware())

	s.router.GET("/parse_site", s.handlers.GetSiteStruct)
	s.router.GET("/details", s.handlers.GetDetails)

	go func() {
		if err := s.router.Start(s.cfg.AppAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Str("function", "Start").Err(err).Msg("Server start error")
		}
	}()

	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.router.Shutdown(ctx)
}
