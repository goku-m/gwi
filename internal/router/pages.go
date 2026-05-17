package router

import (
	"github.com/goku-m/gwi/internal/handler"
	"github.com/labstack/echo/v4"
)

func registerPageRoutes(r *echo.Echo, h *handler.Handlers) {
	r.GET("/", h.Page.Home)
	r.GET("/logs", h.Page.Logs)
}
