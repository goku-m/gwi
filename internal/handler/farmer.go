package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/goku-m/gwi/internal/model"
	"github.com/goku-m/gwi/internal/model/farmer"
	"github.com/goku-m/gwi/internal/server"
	"github.com/goku-m/gwi/internal/service"
	"github.com/labstack/echo/v4"
)

type FarmerHandler struct {
	Handler
	farmerService *service.FarmerService
}

func NewFarmerHandler(s *server.Server, farmerService *service.FarmerService) *FarmerHandler {
	return &FarmerHandler{
		Handler:       NewHandler(s),
		farmerService: farmerService,
	}
}

// ------------------------------------------------------------
// helpers
// ------------------------------------------------------------

func getZoneName(c echo.Context) string {
	return strings.TrimSpace(c.Param("zoneName"))
}

// func parseFloatPtr(v string) (*float64, error) {
// 	v = strings.TrimSpace(v)
// 	if v == "" {
// 		return nil, nil
// 	}
// 	f, err := strconv.ParseFloat(v, 64)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &f, nil
// }

// ------------------------------------------------------------
// PAGE HANDLERS (zone scoped)
// NOTE: Update your router paths accordingly
// e.g. GET /zones/:zoneName/farmers/page
// ------------------------------------------------------------

// ------------------------------------------------------------
// API HANDLERS (zone scoped)
// Routes:
// POST   /zones/:zoneName/farmers
// GET    /zones/:zoneName/farmers
// GET    /zones/:zoneName/farmers/:id
// PUT    /zones/:zoneName/farmers/:id
// DELETE /zones/:zoneName/farmers/:id
// GET    /zones/:zoneName/farmers/stats   (NEW)
// ------------------------------------------------------------

func (h *FarmerHandler) CreateFarmer(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *farmer.CreateFarmerPayload) (*farmer.Farmer, error) {
			zoneName := getZoneName(c)
			return h.farmerService.CreateFarmer(c, zoneName, payload)
		},
		http.StatusCreated,
		&farmer.CreateFarmerPayload{},
	)(c)
}

func (h *FarmerHandler) GetFarmerByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *farmer.GetFarmerByIDPayload) (*farmer.PopulatedFarmer, error) {
			zoneName := getZoneName(c)
			return h.farmerService.GetFarmerByID(c, zoneName, payload.ID)
		},
		http.StatusOK,
		&farmer.GetFarmerByIDPayload{},
	)(c)
}

func (h *FarmerHandler) GetFarmers(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *farmer.GetFarmersQuery) (*model.PaginatedResponse[farmer.PopulatedFarmer], error) {
			zoneName := getZoneName(c)
			return h.farmerService.GetFarmers(c, zoneName, query)
		},
		http.StatusOK,
		&farmer.GetFarmersQuery{},
	)(c)
}

func (h *FarmerHandler) UpdateFarmer(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *farmer.UpdateFarmerPayload) (*farmer.Farmer, error) {
			zoneName := getZoneName(c)
			return h.farmerService.UpdateFarmer(c, zoneName, payload)
		},
		http.StatusOK,
		&farmer.UpdateFarmerPayload{},
	)(c)
}

func (h *FarmerHandler) DeleteFarmer(c echo.Context) error {
	return HandleNoContent(
		h.Handler,
		func(c echo.Context, payload *farmer.DeleteFarmerPayload) error {
			zoneName := getZoneName(c)
			return h.farmerService.DeleteFarmer(c, zoneName, payload.ID)
		},
		http.StatusNoContent,
		&farmer.DeleteFarmerPayload{},
	)(c)
}

// NEW: Zone stats endpoint
func (h *FarmerHandler) GetZoneStats(c echo.Context) error {
	zoneName := getZoneName(c)
	if zoneName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "zoneName is required")
	}

	stats, err := h.farmerService.GetZoneStats(c, zoneName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}


type PullResponse struct {
	Changes map[string]any `json:"changes"`
	// Watermelon expects ms timestamp
	Timestamp int64 `json:"timestamp"`
}

type PushRequest struct {
	Changes map[string]farmer.TableChanges[farmer.FarmerSyncRecord] `json:"changes"`

	// optional metadata if you send chunking info from client
	LastPulledAt *int64 `json:"lastPulledAt,omitempty"`
	ChunkIndex   *int   `json:"chunkIndex,omitempty"`
	ChunkCount   *int   `json:"chunkCount,omitempty"`
}

// GET /zones/:zoneName/sync/pull?lastPulledAt=0
func (h *FarmerHandler) Pull(c echo.Context) error {
	zoneName := c.Param("zoneName")

	lastStr := c.QueryParam("lastPulledAt")
	if lastStr == "" {
		lastStr = "0"
	}

	lastPulledAt, err := strconv.ParseInt(lastStr, 10, 64)
	if err != nil || lastPulledAt < 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid lastPulledAt")
	}

	result, err := h.farmerService.Pull(c.Request().Context(), zoneName, lastPulledAt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, PullResponse{
		Changes: map[string]any{
			"farmers": map[string]any{
				"created": result.Created,
				"updated": result.Updated,
				"deleted": result.Deleted,
			},
		},
		Timestamp: result.Timestamp,
	})
}

// POST /zones/:zoneName/sync/push
func (h *FarmerHandler) Push(c echo.Context) error {
	zoneName := c.Param("zoneName")

	var req PushRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload")
	}

	if err := h.farmerService.Push(c.Request().Context(), zoneName, req.Changes); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"success": true})
}
