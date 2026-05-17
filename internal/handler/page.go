package handler

import (
	"net/http"

	"github.com/goku-m/gwi/internal/pages"
	"github.com/goku-m/gwi/internal/server"
	"github.com/labstack/echo/v4"
)

type PageHandler struct {
	Handler
}

func NewPageHandler(s *server.Server) *PageHandler {
	return &PageHandler{Handler: NewHandler(s)}
}

func (h *PageHandler) Home(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return pages.Home().Render(c.Request().Context(), c.Response().Writer)
}

func (h *PageHandler) Logs(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	return pages.Logs().Render(c.Request().Context(), c.Response().Writer)
}
