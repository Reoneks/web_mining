package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RecoverMiddleware() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			// log.Alert(errb.NewError("RecoverMiddleware").WithMessageAndError("Server fatal message", err))
			return nil
		},
	})
}
