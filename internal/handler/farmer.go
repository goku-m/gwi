package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/goku-m/gwi/internal/model/farmer"
	"github.com/goku-m/gwi/internal/server"
	"github.com/goku-m/gwi/internal/service"
	"github.com/goku-m/gwi/internal/validation"
	"github.com/labstack/echo/v4"
)

type FarmerHandler struct {
	Handler
	farmerService *service.FarmerService
}

var newFarmersStartDate = time.Date(2026, time.May, 6, 0, 0, 0, 0, time.UTC)

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

func getCommunityName(c echo.Context) string {
	return strings.TrimSpace(c.Param("communityName"))
}

func parseStatsDateRange(c echo.Context) (*time.Time, *time.Time, error) {
	const layout = "2006-01-02"

	fromRaw := strings.TrimSpace(c.QueryParam("from"))
	toRaw := strings.TrimSpace(c.QueryParam("to"))

	var fromDate *time.Time
	var toDate *time.Time

	if fromRaw != "" {
		parsed, err := time.Parse(layout, fromRaw)
		if err != nil {
			return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "invalid from date, use YYYY-MM-DD")
		}
		fromDate = &parsed
	}

	if toRaw != "" {
		parsed, err := time.Parse(layout, toRaw)
		if err != nil {
			return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "invalid to date, use YYYY-MM-DD")
		}
		toDate = &parsed
	}

	if fromDate != nil && toDate != nil && fromDate.After(*toDate) {
		return nil, nil, echo.NewHTTPError(http.StatusBadRequest, "from date cannot be after to date")
	}

	return fromDate, toDate, nil
}

func parseLogDate(c echo.Context) (time.Time, error) {
	const layout = "2006-01-02"

	dateRaw := strings.TrimSpace(c.QueryParam("date"))
	if dateRaw == "" {
		return time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "date is required, use YYYY-MM-DD")
	}

	parsed, err := time.Parse(layout, dateRaw)
	if err != nil {
		return time.Time{}, echo.NewHTTPError(http.StatusBadRequest, "invalid date, use YYYY-MM-DD")
	}

	return parsed, nil
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
	query := &farmer.GetFarmersQuery{}
	if err := validation.BindAndValidate(c, query); err != nil {
		return err
	}

	zoneName := getZoneName(c)
	result, err := h.farmerService.GetFarmers(c, zoneName, query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *FarmerHandler) GetEditStatus(c echo.Context) error {
	query := &farmer.GetEditQuery{}

	zoneName := getZoneName(c)
	result, err := h.farmerService.GetEditStatus(c, zoneName, query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
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

	fromDate, toDate, dateErr := parseStatsDateRange(c)
	if dateErr != nil {
		return dateErr
	}

	stats, err := h.farmerService.GetZoneStats(c, zoneName, fromDate, toDate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *FarmerHandler) GetGeneralStats(c echo.Context) error {
	fromDate, toDate, dateErr := parseStatsDateRange(c)
	if dateErr != nil {
		return dateErr
	}

	stats, err := h.farmerService.GetGeneralStats(c, fromDate, toDate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *FarmerHandler) GetCommunityStats(c echo.Context) error {
	zoneName := getZoneName(c)
	if zoneName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "zoneName is required")
	}

	communityName := getCommunityName(c)
	if communityName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "communityName is required")
	}

	fromDate, toDate, dateErr := parseStatsDateRange(c)
	if dateErr != nil {
		return dateErr
	}

	stats, err := h.farmerService.GetCommunityStats(c, zoneName, communityName, fromDate, toDate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *FarmerHandler) GetZoneCommunities(c echo.Context) error {
	zoneName := getZoneName(c)
	if zoneName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "zoneName is required")
	}

	result, err := h.farmerService.GetZoneCommunities(c, zoneName)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, result)
}

func (h *FarmerHandler) GetGeneralNewFarmers(c echo.Context) error {
	stats, err := h.farmerService.GetGeneralNewFarmersCount(c, newFarmersStartDate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *FarmerHandler) GetDailyLogNames(c echo.Context) error {
	logDate, err := parseLogDate(c)
	if err != nil {
		return err
	}

	result, svcErr := h.farmerService.GetDailyLogs(c, logDate)
	if svcErr != nil {
		return svcErr
	}

	return c.JSON(http.StatusOK, result)
}

func (h *FarmerHandler) GetZoneNewFarmers(c echo.Context) error {
	zoneName := getZoneName(c)
	if zoneName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "zoneName is required")
	}

	stats, err := h.farmerService.GetZoneNewFarmersCount(c, zoneName, newFarmersStartDate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *FarmerHandler) GetCommunityNewFarmers(c echo.Context) error {
	zoneName := getZoneName(c)
	if zoneName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "zoneName is required")
	}

	communityName := getCommunityName(c)
	if communityName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "communityName is required")
	}

	stats, err := h.farmerService.GetCommunityNewFarmersCount(c, zoneName, communityName, newFarmersStartDate)
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

	// Watermelon sync expects array fields, not null.
	created := result.Created
	updated := result.Updated
	deleted := result.Deleted
	if created == nil {
		created = []farmer.FarmerSyncRecord{}
	}
	if updated == nil {
		updated = []farmer.FarmerSyncRecord{}
	}
	if deleted == nil {
		deleted = []string{}
	}

	return c.JSON(http.StatusOK, PullResponse{
		Changes: map[string]any{
			"farmers": map[string]any{
				"created": created,
				"updated": updated,
				"deleted": deleted,
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
