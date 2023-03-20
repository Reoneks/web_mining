package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

func RecoverMiddleware() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Error().Str("function", "RecoverMiddleware").Err(err).Bytes("stack", stack).Msg("Server fatal error")
			return nil
		},
	})
}
