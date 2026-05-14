package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/goku-m/gwi/internal/middleware"
	"github.com/goku-m/gwi/internal/model"
	"github.com/goku-m/gwi/internal/model/farmer"
	"github.com/goku-m/gwi/internal/repository"
	"github.com/goku-m/gwi/internal/server"
)

const maxSyncCreatedFarmers = 200

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

func (s *FarmerService) GetEditStatus(ctx echo.Context, zoneName string, query *farmer.GetEditQuery) (*farmer.EditStatus, error) {
	logger := middleware.GetLogger(ctx)

	result, err := s.farmerRepo.GetEditStatus(ctx.Request().Context(), zoneName, query)
	if err != nil {
		logger.Error().Err(err).Str("zone", zoneName).Msg("failed to fetch edit status")
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

func (s *FarmerService) GetZoneStats(ctx echo.Context, zoneName string, fromDate, toDate *time.Time) (*farmer.FarmerStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.farmerRepo.GetZoneStats(ctx.Request().Context(), zoneName, fromDate, toDate)
	if err != nil {
		logger.Error().Err(err).Str("zone", zoneName).Msg("failed to fetch zone farmer statistics")
		return nil, err
	}

	return stats, nil
}

func (s *FarmerService) GetGeneralStats(ctx echo.Context, fromDate, toDate *time.Time) (*farmer.FarmerStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.farmerRepo.GetGeneralStats(ctx.Request().Context(), fromDate, toDate)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch general farmer statistics")
		return nil, err
	}

	return stats, nil
}

func (s *FarmerService) GetCommunityStats(ctx echo.Context, zoneName, communityName string, fromDate, toDate *time.Time) (*farmer.CommunityFarmerStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.farmerRepo.GetCommunityStats(ctx.Request().Context(), zoneName, communityName, fromDate, toDate)
	if err != nil {
		logger.Error().Err(err).
			Str("zone", zoneName).
			Str("community", communityName).
			Msg("failed to fetch community farmer statistics")
		return nil, err
	}

	return stats, nil
}

func (s *FarmerService) GetZoneCommunities(ctx echo.Context, zoneName string) (*farmer.ZoneCommunitiesResponse, error) {
	logger := middleware.GetLogger(ctx)

	communities, err := s.farmerRepo.GetZoneCommunities(ctx.Request().Context(), zoneName)
	if err != nil {
		logger.Error().Err(err).Str("zone", zoneName).Msg("failed to fetch zone communities")
		return nil, err
	}

	return &farmer.ZoneCommunitiesResponse{
		ZoneName:    zoneName,
		Communities: communities,
	}, nil
}

func (s *FarmerService) GetGeneralNewFarmersCount(ctx echo.Context, sinceDate time.Time) (*farmer.NewFarmersStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.farmerRepo.GetGeneralNewFarmersCount(ctx.Request().Context(), sinceDate)
	if err != nil {
		logger.Error().Err(err).Str("since_date", sinceDate.Format("2006-01-02")).Msg("failed to fetch general new farmers count")
		return nil, err
	}

	return stats, nil
}

func (s *FarmerService) GetZoneNewFarmersCount(ctx echo.Context, zoneName string, sinceDate time.Time) (*farmer.NewFarmersStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.farmerRepo.GetZoneNewFarmersCount(ctx.Request().Context(), zoneName, sinceDate)
	if err != nil {
		logger.Error().Err(err).
			Str("zone", zoneName).
			Str("since_date", sinceDate.Format("2006-01-02")).
			Msg("failed to fetch zone new farmers count")
		return nil, err
	}

	return stats, nil
}

func (s *FarmerService) GetCommunityNewFarmersCount(ctx echo.Context, zoneName, communityName string, sinceDate time.Time) (*farmer.NewFarmersStats, error) {
	logger := middleware.GetLogger(ctx)

	stats, err := s.farmerRepo.GetCommunityNewFarmersCount(ctx.Request().Context(), zoneName, communityName, sinceDate)
	if err != nil {
		logger.Error().Err(err).
			Str("zone", zoneName).
			Str("community", communityName).
			Str("since_date", sinceDate.Format("2006-01-02")).
			Msg("failed to fetch community new farmers count")
		return nil, err
	}

	return stats, nil
}

func (s *FarmerService) Pull(ctx context.Context, zoneName string, lastPulledAtMs int64) (*farmer.PullResult, error) {
	created, updated, deleted, ts, err := s.farmerRepo.PullFarmers(ctx, zoneName, lastPulledAtMs)
	if err != nil {
		return nil, err
	}

	if len(created)+len(updated)+len(deleted) > 0 {
		if incErr := s.farmerRepo.IncrementDailySync(ctx, zoneName); incErr != nil {
			s.server.Logger.Warn().
				Err(incErr).
				Str("zone", zoneName).
				Msg("failed to increment daily sync counter after pull")
		}
	}

	return &farmer.PullResult{
		Created:   created,
		Updated:   updated,
		Deleted:   deleted,
		Timestamp: ts,
	}, nil
}

func (s *FarmerService) Push(ctx context.Context, zoneName string, changes map[string]farmer.TableChanges[farmer.FarmerSyncRecord]) error {
	manipulated := false

	farmers, ok := changes["farmers"]
	if ok {
		if len(farmers.Created) > maxSyncCreatedFarmers {
			return fmt.Errorf("sync blocked: created farmers exceeds limit (max %d)", maxSyncCreatedFarmers)
		}

		// Apply created + updated via upsert
		upserts := make([]farmer.FarmerSyncRecord, 0, len(farmers.Created)+len(farmers.Updated))
		upserts = append(upserts, farmers.Created...)
		upserts = append(upserts, farmers.Updated...)
		manipulated = len(upserts)+len(farmers.Deleted) > 0

		if err := s.farmerRepo.PushFarmers(ctx, zoneName, upserts, farmers.Deleted); err != nil {
			return err
		}
	}

	if manipulated {
		if incErr := s.farmerRepo.IncrementDailySync(ctx, zoneName); incErr != nil {
			s.server.Logger.Warn().
				Err(incErr).
				Str("zone", zoneName).
				Msg("failed to increment daily sync counter after push")
		}
	}

	return nil
}
