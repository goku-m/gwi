package router

import (
	"github.com/goku-m/gwi/internal/handler"
	"github.com/goku-m/gwi/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerFarmerRoutes(r *echo.Group, h *handler.FarmerHandler, auth *middleware.AuthMiddleware) {
	// General stats across all zones
	r.GET("/farmers/stats", h.GetGeneralStats)

	// ------------------------------------------------------------
	// Zone-scoped routes
	// Base: /zones/:zoneName
	// ------------------------------------------------------------
	zones := r.Group("/zones/:zoneName")

	// ------------------------------------------------------------
	// Farmer operations (per zone)
	// Base: /zones/:zoneName/farmers
	// ------------------------------------------------------------
	farmers := zones.Group("/farmers")
	//farmers.Use(auth.RequireAuthIP)

	// Collection operations
	farmers.POST("", h.CreateFarmer)
	farmers.GET("", h.GetFarmers)

	// Stats (NEW)
	farmers.GET("/stats", h.GetZoneStats)

	// Community stats
	r.GET("/zones/:zoneName/:communityName/farmers/stats", h.GetCommunityStats)
	r.GET("/zones/:zoneName/communities", h.GetZoneCommunities)

	sync := zones.Group("/sync")

	sync.GET("/pull", h.Pull)
	sync.POST("/push", h.Push)

	// ------------------------------------------------------------
	// Individual farmer operations
	// Base: /zones/:zoneName/farmers/:id
	// ------------------------------------------------------------
	dynamicFarmer := farmers.Group("/:id")
	dynamicFarmer.GET("", h.GetFarmerByID)
	dynamicFarmer.PUT("", h.UpdateFarmer)
	dynamicFarmer.DELETE("", h.DeleteFarmer)
}
