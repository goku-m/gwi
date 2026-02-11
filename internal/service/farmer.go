package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/goku-m/gwi/internal/middleware"
	"github.com/goku-m/gwi/internal/model"
	"github.com/goku-m/gwi/internal/model/farmer"
	"github.com/goku-m/gwi/internal/repository"
	"github.com/goku-m/gwi/internal/server"
)

type FarmerService struct {
	server     *server.Server
	farmerRepo *repository.FarmerRepository
}

func NewFarmerService(server *server.Server, farmerRepo *repository.FarmerRepository) *FarmerService {
	return &FarmerService{
		server:     server,
		farmerRepo: farmerRepo,
	}
}

// ------------------------------------------------------------
// Create (zone scoped)
// ------------------------------------------------------------

func (s *FarmerService) CreateFarmer(ctx echo.Context, zoneName string, payload *farmer.CreateFarmerPayload) (*farmer.Farmer, error) {
	logger := middleware.GetLogger(ctx)

	farmerItem, err := s.farmerRepo.CreateFarmer(ctx.Request().Context(), zoneName, payload)
	if err != nil {
		logger.Error().Err(err).Str("zone", zoneName).Msg("failed to create farmer")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "farmer_created").
		Str("zone", zoneName).
		Str("farmer_id", farmerItem.ID.String()).
		Str("name", farmerItem.Name).
		Str("national_id", farmerItem.NationalID).
		Str("community", farmerItem.Community).
		Msg("Farmer created successfully")

	return farmerItem, nil
}

// ------------------------------------------------------------
// Get by ID (zone scoped)
// ------------------------------------------------------------

func (s *FarmerService) GetFarmerByID(ctx echo.Context, zoneName string, farmerID uuid.UUID) (*farmer.PopulatedFarmer, error) {
	logger := middleware.GetLogger(ctx)

	farmerItem, err := s.farmerRepo.GetFarmerByID(ctx.Request().Context(), zoneName, farmerID)
	if err != nil {
		logger.Error().Err(err).Str("zone", zoneName).Str("farmer_id", farmerID.String()).Msg("failed to fetch farmer by ID")
		return nil, err
	}

	return farmerItem, nil
}

// ------------------------------------------------------------
// List (zone scoped)
// ------------------------------------------------------------

func (s *FarmerService) GetFarmers(ctx echo.Context, zoneName string, query *farmer.GetFarmersQuery) (*model.PaginatedResponse[farmer.PopulatedFarmer], error) {
	logger := middleware.GetLogger(ctx)

	result, err := s.farmerRepo.GetFarmers(ctx.Request().Context(), zoneName, query)
	if err != nil {
		logger.Error().Err(err).Str("zone", zoneName).Msg("failed to fetch farmers")
		return nil, err
	}

	return result, nil
}

// ------------------------------------------------------------
// Update (zone scoped)
// ------------------------------------------------------------

func (s *FarmerService) UpdateFarmer(ctx echo.Context, zoneName string, payload *farmer.UpdateFarmerPayload) (*farmer.Farmer, error) {
	logger := middleware.GetLogger(ctx)

	updatedFarmer, err := s.farmerRepo.UpdateFarmer(ctx.Request().Context(), zoneName, payload)
	if err != nil {
		logger.Error().Err(err).Str("zone", zoneName).Str("farmer_id", payload.ID.String()).Msg("failed to update farmer")
		return nil, err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "farmer_updated").
		Str("zone", zoneName).
		Str("farmer_id", updatedFarmer.ID.String()).
		Str("name", updatedFarmer.Name).
		Str("national_id", updatedFarmer.NationalID).
		Str("community", updatedFarmer.Community).
		Msg("Farmer updated successfully")

	return updatedFarmer, nil
}

// ------------------------------------------------------------
// Delete (zone scoped)
// ------------------------------------------------------------

func (s *FarmerService) DeleteFarmer(ctx echo.Context, zoneName string, farmerID uuid.UUID) error {
	logger := middleware.GetLogger(ctx)

	err := s.farmerRepo.DeleteFarmer(ctx.Request().Context(), zoneName, farmerID)
	if err != nil {
		logger.Error().Err(err).Str("zone", zoneName).Str("farmer_id", farmerID.String()).Msg("failed to delete farmer")
		return err
	}

	// Business event log
	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "farmer_deleted").
		Str("zone", zoneName).
		Str("farmer_id", farmerID.String()).
		Msg("Farmer deleted successfully")

	return nil
}

// ------------------------------------------------------------
// Zone Stats (NEW)
// total farmers, total kgs, total amount, total prefinance
// ------------------------------------------------------------

func (s *FarmerService) GetZoneStats(ctx echo.Context, zoneName string) (*farmer.FarmerStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.farmerRepo.GetZoneStats(ctx.Request().Context(), zoneName)
	if err != nil {
		logger.Error().Err(err).Str("zone", zoneName).Msg("failed to fetch zone farmer statistics")
		return nil, err
	}

	return stats, nil
}

func (s *FarmerService) Pull(ctx context.Context, zoneName string, lastPulledAtMs int64) (*farmer.PullResult, error) {
	created, updated, deleted, ts, err := s.farmerRepo.PullFarmers(ctx, zoneName, lastPulledAtMs)
	if err != nil {
		return nil, err
	}

	return &farmer.PullResult{
		Created:   created,
		Updated:   updated,
		Deleted:   deleted,
		Timestamp: ts,
	}, nil
}

func (s *FarmerService) Push(ctx context.Context, zoneName string, changes map[string]farmer.TableChanges[farmer.FarmerSyncRecord]) error {
	farmers, ok := changes["farmers"]
	if !ok {
		return nil // nothing to do
	}

	// Apply created + updated via upsert
	upserts := make([]farmer.FarmerSyncRecord, 0, len(farmers.Created)+len(farmers.Updated))
	upserts = append(upserts, farmers.Created...)
	upserts = append(upserts, farmers.Updated...)

	return s.farmerRepo.PushFarmers(ctx, zoneName, upserts, farmers.Deleted)
}
