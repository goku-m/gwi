package handler

import (
	"github.com/goku-m/gwi/internal/server"
	"github.com/goku-m/gwi/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
	Farmer  *FarmerHandler
	Auth    *AuthHandler
	Page    *PageHandler
}

func NewHandlers(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
		Farmer:  NewFarmerHandler(s, services.Farmer),
		Auth:    NewAuthHandler(s),
		Page:    NewPageHandler(s),
	}
}
