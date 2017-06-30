package middleware

import (
	"github.com/labstack/echo"
)

// RequestLogDataMiddleware is a middleware to set request info
func RequestLogDataMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoContext echo.Context) error {
			return next(echoContext)
		}
	}
}
