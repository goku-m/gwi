package middleware

import (
	"github.com/labstack/echo/v4"

	"github.com/goku-m/gwi/internal/server"
)

type TracingMiddleware struct {
	server *server.Server
}

func NewTracingMiddleware(s *server.Server) *TracingMiddleware {
	return &TracingMiddleware{server: s}
}

func (tm *TracingMiddleware) RequestTracingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	}
}

func (tm *TracingMiddleware) EnhanceTracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return next
	}
}
