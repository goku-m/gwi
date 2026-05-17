package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goku-m/gwi/internal/errs"
	"github.com/goku-m/gwi/internal/model"
	"github.com/goku-m/gwi/internal/model/farmer"
	"github.com/goku-m/gwi/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FarmerRepository struct {
	server *server.Server
}

func NewFarmerRepository(server *server.Server) *FarmerRepository {
	return &FarmerRepository{server: server}
}

// ------------------------------------------------------------
// Create (zoneName comes from route/middleware, not from payload)
// ------------------------------------------------------------

func (r *FarmerRepository) CreateFarmer(ctx context.Context, zoneName string, payload *farmer.CreateFarmerPayload) (*farmer.Farmer, error) {
	stmt := `
		INSERT INTO farmers (
			zone_name,
			name,
			national_id,
			community,
			prefinance,
			balance,
			total_kg_brought,
			total_amount
		)
		VALUES (
			@zone_name,
			@name,
			@national_id,
			@community,
			COALESCE(@prefinance, 0),
			COALESCE(@balance, 0),
			COALESCE(@total_kg_brought, 0),
			COALESCE(@total_amount, 0)
		)
		RETURNING *
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"zone_name":        zoneName,
		"name":             payload.Name,
		"national_id":      payload.NationalID,
		"community":        payload.Community,
		"prefinance":       payload.Prefinance,
		"balance":          payload.Balance,
		"total_kg_brought": payload.TotalKgBrought,
		"total_amount":     payload.TotalAmount,
	})
	if err != nil {
		return nil, fmt.Errorf("create farmer failed zone=%s national_id=%s: %w", zoneName, payload.NationalID, err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.Farmer])
	if err != nil {
		return nil, fmt.Errorf("collect created farmer failed zone=%s national_id=%s: %w", zoneName, payload.NationalID, err)
	}

	return &item, nil
}

// ------------------------------------------------------------
// Get By ID (scoped to zone)
// ------------------------------------------------------------

func (r *FarmerRepository) GetFarmerByID(ctx context.Context, zoneName string, farmerID uuid.UUID) (*farmer.PopulatedFarmer, error) {
	stmt := `
		SELECT
			t.*
		FROM
			farmers t
		WHERE
			t.id = @id
			AND t.zone_name = @zone_name
			AND t.deleted_at IS NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":        farmerID,
		"zone_name": zoneName,
	})
	if err != nil {
		return nil, fmt.Errorf("get farmer by id failed farmer_id=%s zone=%s: %w", farmerID.String(), zoneName, err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.PopulatedFarmer])
	if err != nil {
		// If you want "not found" specifically, you can map ErrNoRows.
		return nil, fmt.Errorf("collect farmer by id failed farmer_id=%s zone=%s: %w", farmerID.String(), zoneName, err)
	}

	return &item, nil
}

// ------------------------------------------------------------
// Check Exists (scoped to zone)
// ------------------------------------------------------------

func (r *FarmerRepository) CheckFarmerExists(ctx context.Context, zoneName string, farmerID uuid.UUID) (*farmer.Farmer, error) {
	stmt := `
		SELECT
			*
		FROM
			farmers
		WHERE
			id = @id
			AND zone_name = @zone_name
			AND deleted_at IS NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":        farmerID,
		"zone_name": zoneName,
	})
	if err != nil {
		return nil, fmt.Errorf("check farmer exists failed farmer_id=%s zone=%s: %w", farmerID.String(), zoneName, err)
	}

	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.Farmer])
	if err != nil {
		return nil, fmt.Errorf("collect farmer exists failed farmer_id=%s zone=%s: %w", farmerID.String(), zoneName, err)
	}

	return &item, nil
}

// ------------------------------------------------------------
// List (scoped to zone) + paging + safe sort + search + filters
// ------------------------------------------------------------

func (r *FarmerRepository) GetFarmers(
	ctx context.Context,
	zoneName string,
	query *farmer.GetFarmersQuery,
) (*model.PaginatedResponse[farmer.PopulatedFarmer], error) {

	stmt := `
		SELECT
			t.*
		FROM
			farmers t
	`

	args := pgx.NamedArgs{
		"zone_name": zoneName,
	}
	conditions := []string{"t.zone_name = @zone_name", "t.deleted_at IS NULL"}

	if query != nil {
		if query.Search != nil && strings.TrimSpace(*query.Search) != "" {
			conditions = append(conditions, "(t.name ILIKE @search OR t.national_id ILIKE @search)")
			args["search"] = "%" + strings.TrimSpace(*query.Search) + "%"
		}

		if query.Community != nil && strings.TrimSpace(*query.Community) != "" {
			conditions = append(conditions, "LOWER(BTRIM(t.community)) = LOWER(BTRIM(@community))")
			args["community"] = strings.TrimSpace(*query.Community)
		}

		if query.HasDebt != nil {
			if *query.HasDebt {
				conditions = append(conditions, "t.balance > 0")
			} else {
				conditions = append(conditions, "t.balance = 0")
			}
		}
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	// ----- count query -----
	countStmt := "SELECT COUNT(*) FROM farmers t WHERE " + strings.Join(conditions, " AND ")

	var total int
	if err := r.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&total); err != nil {
		return nil, fmt.Errorf("failed to count farmers zone=%s: %w", zoneName, err)
	}

	// ----- pagination defaults -----
	page := 1
	limit := 20
	if query != nil {
		if query.Page != nil && *query.Page > 0 {
			page = *query.Page
		}
		if query.Limit != nil && *query.Limit > 0 {
			limit = *query.Limit
		}
	}

	// ----- safe sorting -----
	sortCol := "created_at"
	orderDesc := true

	allowedSort := map[string]bool{
		"created_at":       true,
		"updated_at":       true,
		"name":             true,
		"community":        true,
		"national_id":      true,
		"balance":          true,
		"total_kg_brought": true,
		"total_amount":     true,
	}

	if query != nil && query.Sort != nil && allowedSort[*query.Sort] {
		sortCol = *query.Sort
	}
	if query != nil && query.Order != nil && strings.EqualFold(*query.Order, "asc") {
		orderDesc = false
	}

	stmt += " ORDER BY t." + sortCol
	if orderDesc {
		stmt += " DESC"
	} else {
		stmt += " ASC"
	}

	// ----- pagination -----
	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = limit
	args["offset"] = (page - 1) * limit

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("get farmers failed zone=%s: %w", zoneName, err)
	}
	defer rows.Close()

	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[farmer.PopulatedFarmer])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[farmer.PopulatedFarmer]{
				Data:       []farmer.PopulatedFarmer{},
				Page:       page,
				Limit:      limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("collect farmers failed zone=%s: %w", zoneName, err)
	}

	return &model.PaginatedResponse[farmer.PopulatedFarmer]{
		Data:       items,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: (total + limit - 1) / limit,
	}, nil
}

// ------------------------------------------------------------
// Update (scoped to zone)
// ------------------------------------------------------------

func (r *FarmerRepository) UpdateFarmer(ctx context.Context, zoneName string, payload *farmer.UpdateFarmerPayload) (*farmer.Farmer, error) {
	stmt := "UPDATE farmers SET "
	args := pgx.NamedArgs{
		"farmer_id": payload.ID,
		"zone_name": zoneName,
	}

	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = strings.TrimSpace(*payload.Name)
	}
	if payload.NationalID != nil {
		setClauses = append(setClauses, "national_id = @national_id")
		args["national_id"] = strings.TrimSpace(*payload.NationalID)
	}
	if payload.Community != nil {
		setClauses = append(setClauses, "community = @community")
		args["community"] = strings.TrimSpace(*payload.Community)
	}

	if payload.Prefinance != nil {
		setClauses = append(setClauses, "prefinance = @prefinance")
		args["prefinance"] = *payload.Prefinance
	}
	if payload.Balance != nil {
		setClauses = append(setClauses, "balance = @balance")
		args["balance"] = *payload.Balance
	}
	if payload.TotalKgBrought != nil {
		setClauses = append(setClauses, "total_kg_brought = @total_kg_brought")
		args["total_kg_brought"] = *payload.TotalKgBrought
	}
	if payload.TotalAmount != nil {
		setClauses = append(setClauses, "total_amount = @total_amount")
		args["total_amount"] = *payload.TotalAmount
	}

	if len(setClauses) == 0 {
		return nil, errs.NewBadRequestError("no fields to update", false, nil, nil, nil)
	}

	// optional: ensure updated_at moves even if your trigger isn't present
	setClauses = append(setClauses, "updated_at = now()")

	stmt += strings.Join(setClauses, ", ")
	stmt += " WHERE id = @farmer_id AND zone_name = @zone_name RETURNING *"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("update farmer failed id=%s zone=%s: %w", payload.ID.String(), zoneName, err)
	}

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.Farmer])
	if err != nil {
		return nil, fmt.Errorf("collect updated farmer failed id=%s zone=%s: %w", payload.ID.String(), zoneName, err)
	}

	return &updated, nil
}

// ------------------------------------------------------------
// Delete (scoped to zone)
// ------------------------------------------------------------

func (r *FarmerRepository) DeleteFarmer(ctx context.Context, zoneName string, farmerID uuid.UUID) error {
	stmt := `
		UPDATE farmers
		SET
			deleted_at = (EXTRACT(EPOCH FROM now()) * 1000)::bigint,
			updated_at = now()
		WHERE
			id = @farmer_id
			AND zone_name = @zone_name
			AND deleted_at IS NULL
	`

	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"farmer_id": farmerID,
		"zone_name": zoneName,
	})
	if err != nil {
		return fmt.Errorf("delete farmer failed id=%s zone=%s: %w", farmerID.String(), zoneName, err)
	}

	if result.RowsAffected() == 0 {
		code := "FARMER_NOT_FOUND"
		return errs.NewNotFoundError("farmer not found", false, &code)
	}

	return nil
}

// ------------------------------------------------------------
// Zone Stats (NEW): totals per zone
// total farmers, total kgs, total amount, total prefinance
// ------------------------------------------------------------

func (r *FarmerRepository) GetZoneStats(ctx context.Context, zoneName string, fromDate, toDate *time.Time) (*farmer.FarmerStats, error) {
	stmt := `
		SELECT
			@zone_name::text AS zone_name,
			COUNT(*) AS total_farmers,
			COUNT(DISTINCT NULLIF(BTRIM(community), '')) AS total_communities,
			COALESCE((
				SELECT zds.sync_count
				FROM zone_daily_syncs zds
				WHERE LOWER(BTRIM(zds.zone_name)) = LOWER(BTRIM(@zone_name))
				  AND zds.sync_date = CURRENT_DATE
			), 0)::int AS daily_syncs,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN total_kg_brought ELSE 0 END), 0)::float8 AS total_kg_brought,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN total_amount ELSE 0 END), 0)::float8 AS total_amount,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN prefinance ELSE 0 END), 0)::float8 AS total_prefinance,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN balance ELSE 0 END), 0)::float8 AS total_balance
		FROM
			farmers
		WHERE
			LOWER(BTRIM(zone_name)) = LOWER(BTRIM(@zone_name))
			AND deleted_at IS NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"zone_name": zoneName,
		"from_date": fromDate,
		"to_date":   toDate,
	})
	if err != nil {
		return nil, fmt.Errorf("get zone stats failed zone=%s: %w", zoneName, err)
	}

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.FarmerStats])
	if err != nil {
		return nil, fmt.Errorf("collect zone stats failed zone=%s: %w", zoneName, err)
	}

	return &stats, nil
}

func (r *FarmerRepository) GetGeneralStats(ctx context.Context, fromDate, toDate *time.Time) (*farmer.FarmerStats, error) {
	stmt := `
		SELECT
			'General'::text AS zone_name,
			COUNT(*) AS total_farmers,
			COUNT(DISTINCT NULLIF(BTRIM(community), '')) AS total_communities,
			COALESCE((
				SELECT SUM(zds.sync_count)
				FROM zone_daily_syncs zds
				WHERE zds.sync_date = CURRENT_DATE
			), 0)::int AS daily_syncs,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN total_kg_brought ELSE 0 END), 0)::float8 AS total_kg_brought,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN total_amount ELSE 0 END), 0)::float8 AS total_amount,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN prefinance ELSE 0 END), 0)::float8 AS total_prefinance,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN balance ELSE 0 END), 0)::float8 AS total_balance
		FROM
			farmers
		WHERE
			deleted_at IS NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"from_date": fromDate,
		"to_date":   toDate,
	})
	if err != nil {
		return nil, fmt.Errorf("get general stats failed: %w", err)
	}

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.FarmerStats])
	if err != nil {
		return nil, fmt.Errorf("collect general stats failed: %w", err)
	}

	return &stats, nil
}

func (r *FarmerRepository) GetCommunityStats(ctx context.Context, zoneName, communityName string, fromDate, toDate *time.Time) (*farmer.CommunityFarmerStats, error) {
	stmt := `
		SELECT
			@zone_name::text AS zone_name,
			@community_name::text AS community_name,
			COUNT(*) AS total_farmers,
			COALESCE((
				SELECT zds.sync_count
				FROM zone_daily_syncs zds
				WHERE LOWER(BTRIM(zds.zone_name)) = LOWER(BTRIM(@zone_name))
				  AND zds.sync_date = CURRENT_DATE
			), 0)::int AS daily_syncs,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN total_kg_brought ELSE 0 END), 0)::float8 AS total_kg_brought,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN total_amount ELSE 0 END), 0)::float8 AS total_amount,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN prefinance ELSE 0 END), 0)::float8 AS total_prefinance,
			COALESCE(SUM(CASE WHEN (@from_date::date IS NULL OR updated_at::date >= @from_date::date) AND (@to_date::date IS NULL OR updated_at::date <= @to_date::date) THEN balance ELSE 0 END), 0)::float8 AS total_balance
		FROM
			farmers
		WHERE
			LOWER(BTRIM(zone_name)) = LOWER(BTRIM(@zone_name))
			AND LOWER(BTRIM(community)) = LOWER(BTRIM(@community_name))
			AND deleted_at IS NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"zone_name":      zoneName,
		"community_name": communityName,
		"from_date":      fromDate,
		"to_date":        toDate,
	})
	if err != nil {
		return nil, fmt.Errorf("get community stats failed zone=%s community=%s: %w", zoneName, communityName, err)
	}

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.CommunityFarmerStats])
	if err != nil {
		return nil, fmt.Errorf("collect community stats failed zone=%s community=%s: %w", zoneName, communityName, err)
	}

	return &stats, nil
}

func (r *FarmerRepository) IncrementDailySync(ctx context.Context, zoneName string) error {
	stmt := `
		INSERT INTO zone_daily_syncs (zone_name, sync_date, sync_count, created_at, updated_at)
		VALUES (@zone_name, CURRENT_DATE, 1, NOW(), NOW())
		ON CONFLICT (zone_name, sync_date)
		DO UPDATE SET
			sync_count = zone_daily_syncs.sync_count + 1,
			updated_at = NOW()
	`

	_, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"zone_name": strings.TrimSpace(zoneName),
	})
	if err != nil {
		return fmt.Errorf("increment daily sync failed zone=%s: %w", zoneName, err)
	}

	return nil
}

func (r *FarmerRepository) GetZoneCommunities(ctx context.Context, zoneName string) ([]string, error) {
	stmt := `
		SELECT DISTINCT BTRIM(community) AS community_name
		FROM farmers
		WHERE LOWER(BTRIM(zone_name)) = LOWER(BTRIM(@zone_name))
		  AND deleted_at IS NULL
		  AND BTRIM(community) <> ''
		ORDER BY community_name ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"zone_name": zoneName,
	})
	if err != nil {
		return nil, fmt.Errorf("get zone communities failed zone=%s: %w", zoneName, err)
	}
	defer rows.Close()

	communities := make([]string, 0)
	for rows.Next() {
		var community string
		if err := rows.Scan(&community); err != nil {
			return nil, fmt.Errorf("scan zone community failed zone=%s: %w", zoneName, err)
		}
		communities = append(communities, community)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate zone communities failed zone=%s: %w", zoneName, err)
	}

	return communities, nil
}

func (r *FarmerRepository) GetGeneralNewFarmersCount(ctx context.Context, sinceDate time.Time) (*farmer.NewFarmersStats, error) {
	stmt := `
		SELECT
			COUNT(*)::int AS new_farmers,
			@since_date::date::text AS since_date
		FROM farmers
		WHERE created_at::date >= @since_date::date
		  AND deleted_at IS NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"since_date": sinceDate,
	})
	if err != nil {
		return nil, fmt.Errorf("get general new farmers failed since=%s: %w", sinceDate.Format("2006-01-02"), err)
	}

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.NewFarmersStats])
	if err != nil {
		return nil, fmt.Errorf("collect general new farmers failed since=%s: %w", sinceDate.Format("2006-01-02"), err)
	}

	return &stats, nil
}

func (r *FarmerRepository) GetDailyLogs(ctx context.Context, logDate time.Time) ([]farmer.DailyLogEntry, error) {
	stmt := `
		WITH created_logs AS (
			SELECT
				MIN(created_at) AS event_at,
				TO_CHAR(MIN(created_at), 'DD/MM/YY') AS date_label,
				TO_CHAR(MIN(created_at), 'HH24:MI') AS time_label,
				BTRIM(created_by) AS created_by,
				''::text AS updated_by,
				BTRIM(zone_name) AS zone_name,
				COUNT(*)::int AS farmers_count,
				COUNT(DISTINCT NULLIF(BTRIM(community), ''))::int AS communities_count,
				'created'::text AS action,
				0::float8 AS amount,
				0::float8 AS weight_kg
			FROM farmers
			WHERE deleted_at IS NULL
			  AND created_at::date = @log_date::date
			  AND BTRIM(created_by) <> ''
			  AND LENGTH(BTRIM(created_by)) > 1
			  AND UPPER(BTRIM(created_by)) NOT IN ('NIL', 'UNKNOWN')
			GROUP BY BTRIM(created_by), BTRIM(zone_name)
		),
		updated_logs AS (
			SELECT
				MAX(updated_at) AS event_at,
				TO_CHAR(MAX(updated_at), 'DD/MM/YY') AS date_label,
				TO_CHAR(MAX(updated_at), 'HH24:MI') AS time_label,
				''::text AS created_by,
				BTRIM(updated_by) AS updated_by,
				BTRIM(zone_name) AS zone_name,
				0::int AS farmers_count,
				COUNT(DISTINCT NULLIF(BTRIM(community), ''))::int AS communities_count,
				'updated'::text AS action,
				COALESCE(SUM(prefinance), 0)::float8 AS amount,
				0::float8 AS weight_kg
			FROM farmers
			WHERE deleted_at IS NULL
			  AND updated_at::date = @log_date::date
			  AND BTRIM(updated_by) <> ''
			  AND LENGTH(BTRIM(updated_by)) > 1
			  AND UPPER(BTRIM(updated_by)) NOT IN ('NIL', 'UNKNOWN')
			GROUP BY BTRIM(updated_by), BTRIM(zone_name)
		),
		weighed_logs AS (
			SELECT
				MAX(updated_at) AS event_at,
				TO_CHAR(MAX(updated_at), 'DD/MM/YY') AS date_label,
				TO_CHAR(MAX(updated_at), 'HH24:MI') AS time_label,
				''::text AS created_by,
				BTRIM(updated_by) AS updated_by,
				BTRIM(zone_name) AS zone_name,
				0::int AS farmers_count,
				COUNT(DISTINCT NULLIF(BTRIM(community), ''))::int AS communities_count,
				'weighed'::text AS action,
				COALESCE(SUM(total_amount), 0)::float8 AS amount,
				COALESCE(SUM(total_kg_brought), 0)::float8 AS weight_kg
			FROM farmers
			WHERE deleted_at IS NULL
			  AND updated_at::date = @log_date::date
			  AND BTRIM(updated_by) <> ''
			  AND LENGTH(BTRIM(updated_by)) > 1
			  AND UPPER(BTRIM(updated_by)) NOT IN ('NIL', 'UNKNOWN')
			GROUP BY BTRIM(updated_by), BTRIM(zone_name)
			HAVING COALESCE(SUM(total_kg_brought), 0) > 0
		)
		SELECT
			date_label,
			time_label,
			created_by,
			updated_by,
			zone_name,
			farmers_count,
			communities_count,
			action,
			amount,
			weight_kg
		FROM (
			SELECT * FROM created_logs
			UNION ALL
			SELECT * FROM updated_logs
			UNION ALL
			SELECT * FROM weighed_logs
		) logs
		ORDER BY event_at ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"log_date": logDate,
	})
	if err != nil {
		return nil, fmt.Errorf("get daily logs failed date=%s: %w", logDate.Format("2006-01-02"), err)
	}
	defer rows.Close()

	logs := make([]farmer.DailyLogEntry, 0)
	for rows.Next() {
		var entry farmer.DailyLogEntry
		if err := rows.Scan(
			&entry.Date,
			&entry.Time,
			&entry.CreatedBy,
			&entry.UpdatedBy,
			&entry.ZoneName,
			&entry.Count,
			&entry.Communities,
			&entry.Action,
			&entry.Amount,
			&entry.WeightKg,
		); err != nil {
			return nil, fmt.Errorf("scan daily log failed date=%s: %w", logDate.Format("2006-01-02"), err)
		}
		logs = append(logs, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate daily logs failed date=%s: %w", logDate.Format("2006-01-02"), err)
	}

	return logs, nil
}

func (r *FarmerRepository) GetDailyLogNames(ctx context.Context, logDate time.Time) ([]string, error) {
	stmt := `
		WITH names AS (
			SELECT BTRIM(created_by) AS name
			FROM farmers
			WHERE deleted_at IS NULL
			  AND created_at::date = @log_date::date
			  AND BTRIM(created_by) <> ''
			  AND LENGTH(BTRIM(created_by)) > 1
			  AND UPPER(BTRIM(created_by)) NOT IN ('NIL', 'UNKNOWN')
			UNION
			SELECT BTRIM(updated_by) AS name
			FROM farmers
			WHERE deleted_at IS NULL
			  AND updated_at::date = @log_date::date
			  AND BTRIM(updated_by) <> ''
			  AND LENGTH(BTRIM(updated_by)) > 1
			  AND UPPER(BTRIM(updated_by)) NOT IN ('NIL', 'UNKNOWN')
		)
		SELECT name
		FROM names
		ORDER BY LOWER(name) ASC
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"log_date": logDate,
	})
	if err != nil {
		return nil, fmt.Errorf("get daily log names failed date=%s: %w", logDate.Format("2006-01-02"), err)
	}
	defer rows.Close()

	names := make([]string, 0)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("scan daily log name failed date=%s: %w", logDate.Format("2006-01-02"), err)
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate daily log names failed date=%s: %w", logDate.Format("2006-01-02"), err)
	}

	return names, nil
}

func (r *FarmerRepository) GetZoneNewFarmersCount(ctx context.Context, zoneName string, sinceDate time.Time) (*farmer.NewFarmersStats, error) {
	stmt := `
		SELECT
			COUNT(*)::int AS new_farmers,
			@since_date::date::text AS since_date
		FROM farmers
		WHERE LOWER(BTRIM(zone_name)) = LOWER(BTRIM(@zone_name))
		  AND created_at::date >= @since_date::date
		  AND deleted_at IS NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"zone_name":  zoneName,
		"since_date": sinceDate,
	})
	if err != nil {
		return nil, fmt.Errorf("get zone new farmers failed zone=%s since=%s: %w", zoneName, sinceDate.Format("2006-01-02"), err)
	}

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.NewFarmersStats])
	if err != nil {
		return nil, fmt.Errorf("collect zone new farmers failed zone=%s since=%s: %w", zoneName, sinceDate.Format("2006-01-02"), err)
	}

	return &stats, nil
}

func (r *FarmerRepository) GetCommunityNewFarmersCount(ctx context.Context, zoneName, communityName string, sinceDate time.Time) (*farmer.NewFarmersStats, error) {
	stmt := `
		SELECT
			COUNT(*)::int AS new_farmers,
			@since_date::date::text AS since_date
		FROM farmers
		WHERE LOWER(BTRIM(zone_name)) = LOWER(BTRIM(@zone_name))
		  AND LOWER(BTRIM(community)) = LOWER(BTRIM(@community_name))
		  AND created_at::date >= @since_date::date
		  AND deleted_at IS NULL
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"zone_name":      zoneName,
		"community_name": communityName,
		"since_date":     sinceDate,
	})
	if err != nil {
		return nil, fmt.Errorf("get community new farmers failed zone=%s community=%s since=%s: %w", zoneName, communityName, sinceDate.Format("2006-01-02"), err)
	}

	stats, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[farmer.NewFarmersStats])
	if err != nil {
		return nil, fmt.Errorf("collect community new farmers failed zone=%s community=%s since=%s: %w", zoneName, communityName, sinceDate.Format("2006-01-02"), err)
	}

	return &stats, nil
}

func (r *FarmerRepository) GetEditStatus(
	ctx context.Context,
	zoneName string,
	query *farmer.GetEditQuery,
) (*farmer.EditStatus, error) {
	_ = zoneName
	_ = query

	stmt := `
		SELECT should_edit
		FROM edit_status
		ORDER BY updated_at DESC
		LIMIT 1
	`

	var status farmer.EditStatus
	if err := r.server.DB.Pool.QueryRow(ctx, stmt).Scan(&status.ShouldEdit); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Safe default when no status row exists yet.
			status.ShouldEdit = true
			return &status, nil
		}
		return nil, fmt.Errorf("get edit status failed: %w", err)
	}

	return &status, nil
}

type SyncFarmerRow struct {
	ID            string
	ZoneName      string
	Name          string
	NationalID    string
	Community     string
	Prefinance    float64
	Balance       float64
	TotalKg       float64
	TotalAmount   float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CreatedBy     string
	UpdatedBy     string
	DeletedAtNull bool
	DeletedAtMs   int64
}

type SyncRepository struct {
	pool *pgxpool.Pool
}

func NewSyncRepository(pool *pgxpool.Pool) *SyncRepository {
	return &SyncRepository{pool: pool}
}

func msToTime(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond)).UTC()
}
func timeToMs(t time.Time) int64 {
	return t.UTC().UnixNano() / int64(time.Millisecond)
}

func normalizeCreatedBy(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "UNKNOWN"
	}
	return v
}

func normalizeUpdatedBy(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "NIL"
	}
	return v
}

func (r *FarmerRepository) PullFarmers(ctx context.Context, zoneName string, lastPulledAtMs int64) (created, updated []farmer.FarmerSyncRecord, deletedIDs []string, timestampMs int64, err error) {
	now := time.Now().UTC()
	timestampMs = timeToMs(now)

	// First sync = all data as created
	if lastPulledAtMs == 0 {
		stmt := `
			SELECT
				id::text,
				zone_name,
				name,
				national_id,
				community,
				prefinance::float8,
				balance::float8,
				total_kg_brought::float8,
				total_amount::float8,
				created_at,
				updated_at,
				created_by,
				updated_by,
				(deleted_at IS NULL) AS deleted_at_null,
				COALESCE(deleted_at, 0) AS deleted_at_ms
			FROM farmers
			WHERE zone_name = @zone_name
			  AND deleted_at IS NULL
		`
		rows, qerr := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{"zone_name": zoneName})
		if qerr != nil {
			return nil, nil, nil, 0, fmt.Errorf("pull farmers first sync: %w", qerr)
		}
		defer rows.Close()

		for rows.Next() {
			var row SyncFarmerRow
			if err := rows.Scan(
				&row.ID, &row.ZoneName, &row.Name, &row.NationalID, &row.Community,
				&row.Prefinance, &row.Balance, &row.TotalKg, &row.TotalAmount,
				&row.CreatedAt, &row.UpdatedAt, &row.CreatedBy, &row.UpdatedBy, &row.DeletedAtNull, &row.DeletedAtMs,
			); err != nil {
				return nil, nil, nil, 0, fmt.Errorf("scan farmer row: %w", err)
			}
			created = append(created, farmerSyncFromRow(row))
		}
		return created, nil, nil, timestampMs, nil
	}

	last := msToTime(lastPulledAtMs)

	// created: created_at > last, not deleted
	createdStmt := `
		SELECT
			id::text,
			zone_name,
			name,
			national_id,
			community,
			prefinance::float8,
			balance::float8,
			total_kg_brought::float8,
			total_amount::float8,
			created_at,
			updated_at,
			created_by,
			updated_by,
			(deleted_at IS NULL) AS deleted_at_null,
			COALESCE(deleted_at, 0) AS deleted_at_ms
		FROM farmers
		WHERE zone_name = @zone_name
		  AND deleted_at IS NULL
		  AND created_at > @last
	`
	// updated: updated_at > last, created_at <= last, not deleted
	updatedStmt := `
		SELECT
			id::text,
			zone_name,
			name,
			national_id,
			community,
			prefinance::float8,
			balance::float8,
			total_kg_brought::float8,
			total_amount::float8,
			created_at,
			updated_at,
			created_by,
			updated_by,
			(deleted_at IS NULL) AS deleted_at_null,
			COALESCE(deleted_at, 0) AS deleted_at_ms
		FROM farmers
		WHERE zone_name = @zone_name
		  AND deleted_at IS NULL
		  AND updated_at > @last
		  AND created_at <= @last
	`
	// deleted: deleted_at > last
	deletedStmt := `
		SELECT DISTINCT id
		FROM (
			SELECT id::text AS id
			FROM farmers
			WHERE zone_name = @zone_name
			  AND deleted_at IS NOT NULL
			  AND deleted_at > @last_ms
			UNION ALL
			SELECT farmer_id::text AS id
			FROM farmer_deletions
			WHERE zone_name = @zone_name
			  AND deleted_at > @last_ms
		) d
	`

	args := pgx.NamedArgs{"zone_name": zoneName, "last": last, "last_ms": lastPulledAtMs}

	// created
	{
		rows, qerr := r.server.DB.Pool.Query(ctx, createdStmt, args)
		if qerr != nil {
			return nil, nil, nil, 0, fmt.Errorf("pull created farmers: %w", qerr)
		}
		defer rows.Close()
		for rows.Next() {
			var row SyncFarmerRow
			if err := rows.Scan(
				&row.ID, &row.ZoneName, &row.Name, &row.NationalID, &row.Community,
				&row.Prefinance, &row.Balance, &row.TotalKg, &row.TotalAmount,
				&row.CreatedAt, &row.UpdatedAt, &row.CreatedBy, &row.UpdatedBy, &row.DeletedAtNull, &row.DeletedAtMs,
			); err != nil {
				return nil, nil, nil, 0, fmt.Errorf("scan created farmer: %w", err)
			}
			created = append(created, farmerSyncFromRow(row))
		}
	}

	// updated
	{
		rows, qerr := r.server.DB.Pool.Query(ctx, updatedStmt, args)
		if qerr != nil {
			return nil, nil, nil, 0, fmt.Errorf("pull updated farmers: %w", qerr)
		}
		defer rows.Close()
		for rows.Next() {
			var row SyncFarmerRow
			if err := rows.Scan(
				&row.ID, &row.ZoneName, &row.Name, &row.NationalID, &row.Community,
				&row.Prefinance, &row.Balance, &row.TotalKg, &row.TotalAmount,
				&row.CreatedAt, &row.UpdatedAt, &row.CreatedBy, &row.UpdatedBy, &row.DeletedAtNull, &row.DeletedAtMs,
			); err != nil {
				return nil, nil, nil, 0, fmt.Errorf("scan updated farmer: %w", err)
			}
			updated = append(updated, farmerSyncFromRow(row))
		}
	}

	// deleted ids
	{
		rows, qerr := r.server.DB.Pool.Query(ctx, deletedStmt, args)
		if qerr != nil {
			return nil, nil, nil, 0, fmt.Errorf("pull deleted farmers: %w", qerr)
		}
		defer rows.Close()
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return nil, nil, nil, 0, fmt.Errorf("scan deleted id: %w", err)
			}
			deletedIDs = append(deletedIDs, id)
		}
	}

	return created, updated, deletedIDs, timestampMs, nil
}

func (r *FarmerRepository) PushFarmers(
	ctx context.Context,
	zoneName string,
	upserts []farmer.FarmerSyncRecord,
	deletedIDs []string,
) error {
	tx, err := r.server.DB.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	upsertSQL := `
		INSERT INTO farmers (
			id, zone_name, name, national_id, community,
			prefinance, balance, total_kg_brought, total_amount,
			created_at, updated_at, created_by, updated_by, deleted_at
		) VALUES (
			@id, @zone_name, @name, @national_id, @community,
			@prefinance, @balance, @total_kg_brought, @total_amount,
			to_timestamp(@created_at_ms::double precision / 1000.0),
			to_timestamp(@updated_at_ms::double precision / 1000.0),
			@created_by,
			@updated_by,
			NULL
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			national_id = EXCLUDED.national_id,
			community = EXCLUDED.community,
			prefinance = EXCLUDED.prefinance,
			balance = EXCLUDED.balance,
			total_kg_brought = EXCLUDED.total_kg_brought,
			total_amount = EXCLUDED.total_amount,
			updated_at = EXCLUDED.updated_at,
			updated_by = EXCLUDED.updated_by,
			deleted_at = NULL
		WHERE farmers.zone_name = EXCLUDED.zone_name
	`

	deleteSQL := `
		UPDATE farmers
		SET deleted_at = (EXTRACT(EPOCH FROM now()) * 1000)::bigint, updated_at = now()
		WHERE id = @id AND zone_name = @zone_name
	`
	clearDeletionTombstoneSQL := `
		DELETE FROM farmer_deletions
		WHERE farmer_id = @id::uuid AND zone_name = @zone_name
	`
	upsertDeletionTombstoneSQL := `
		INSERT INTO farmer_deletions (farmer_id, zone_name, deleted_at)
		VALUES (@id::uuid, @zone_name, (EXTRACT(EPOCH FROM NOW()) * 1000)::bigint)
		ON CONFLICT (farmer_id, zone_name) DO UPDATE
		SET deleted_at = EXCLUDED.deleted_at
	`

	// Build an in-zone identity index so we can skip duplicates during sync.
	type existingIdentity struct {
		ID         uuid.UUID
		Name       string
		Community  string
		NationalID string
	}

	existingRows, err := r.server.DB.Pool.Query(ctx, `
		SELECT id, name, community, national_id
		FROM farmers
		WHERE LOWER(BTRIM(zone_name)) = LOWER(BTRIM(@zone_name))
		  AND deleted_at IS NULL
	`, pgx.NamedArgs{"zone_name": zoneName})
	if err != nil {
		return fmt.Errorf("load existing farmers for duplicate check failed zone=%s: %w", zoneName, err)
	}
	defer existingRows.Close()

	existingByID := make(map[string]string)
	existingByIdentity := make(map[string]string)
	for existingRows.Next() {
		var row existingIdentity
		if scanErr := existingRows.Scan(&row.ID, &row.Name, &row.Community, &row.NationalID); scanErr != nil {
			return fmt.Errorf("scan existing farmer for duplicate check failed zone=%s: %w", zoneName, scanErr)
		}
		idStr := row.ID.String()
		identityKey := normalizeFarmerIdentityKey(row.Name, row.Community, row.NationalID)
		existingByID[idStr] = identityKey
		if _, present := existingByIdentity[identityKey]; !present {
			existingByIdentity[identityKey] = idStr
		}
	}
	if err := existingRows.Err(); err != nil {
		return fmt.Errorf("iterate existing farmers for duplicate check failed zone=%s: %w", zoneName, err)
	}

	b := &pgx.Batch{}

	for _, f := range upserts {
		identityKey := normalizeFarmerIdentityKey(f.Name, f.Community, f.NationalID)
		existingIdentityKeyForID, existsByID := existingByID[f.ID]
		conflictingID, existsByIdentity := existingByIdentity[identityKey]
		if shouldSkipSyncUpsert(existsByID, existsByIdentity, f.ID, conflictingID) {
			r.server.Logger.Info().
				Str("zone", zoneName).
				Str("incoming_id", f.ID).
				Str("conflicting_id", conflictingID).
				Str("name", f.Name).
				Str("community", f.Community).
				Str("national_id", f.NationalID).
				Msg("sync upsert skipped due to duplicate farmer identity in zone")
			b.Queue(upsertDeletionTombstoneSQL, pgx.NamedArgs{
				"id":        f.ID,
				"zone_name": zoneName,
			})
			continue
		}

		b.Queue(upsertSQL, pgx.NamedArgs{
			"id":               f.ID,
			"zone_name":        zoneName, // enforce route zone
			"name":             f.Name,
			"national_id":      f.NationalID,
			"community":        f.Community,
			"prefinance":       f.Prefinance,
			"balance":          f.Balance,
			"total_kg_brought": f.TotalKgBrought,
			"total_amount":     f.TotalAmount,
			"created_at_ms":    f.CreatedAt,
			"updated_at_ms":    f.UpdatedAt,
			"created_by":       normalizeCreatedBy(f.AddedBy),
			"updated_by":       normalizeUpdatedBy(f.UpdatedBy),
		})
		b.Queue(clearDeletionTombstoneSQL, pgx.NamedArgs{
			"id":        f.ID,
			"zone_name": zoneName,
		})

		if existsByID {
			delete(existingByIdentity, existingIdentityKeyForID)
		}
		existingByID[f.ID] = identityKey
		existingByIdentity[identityKey] = f.ID
	}

	for _, id := range deletedIDs {
		b.Queue(deleteSQL, pgx.NamedArgs{
			"id":        id,
			"zone_name": zoneName,
		})
	}

	br := tx.SendBatch(ctx, b)

	for i := 0; i < b.Len(); i++ {
		if _, err := br.Exec(); err != nil {
			_ = br.Close()
			return fmt.Errorf("batch exec failed: %w", err)
		}
	}

	// Close batch results before commit; otherwise the tx connection remains busy.
	if err := br.Close(); err != nil {
		return fmt.Errorf("close batch: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func farmerSyncFromRow(row SyncFarmerRow) farmer.FarmerSyncRecord {
	var deletedAt *int64
	if !row.DeletedAtNull {
		deletedAt = &row.DeletedAtMs
	}

	return farmer.FarmerSyncRecord{
		ID:             row.ID,
		ZoneName:       row.ZoneName,
		Name:           row.Name,
		NationalID:     row.NationalID,
		Community:      row.Community,
		Prefinance:     row.Prefinance,
		Balance:        row.Balance,
		TotalKgBrought: row.TotalKg,
		TotalAmount:    row.TotalAmount,
		CreatedAt:      timeToMs(row.CreatedAt),
		UpdatedAt:      timeToMs(row.UpdatedAt),
		AddedBy:        row.CreatedBy,
		UpdatedBy:      row.UpdatedBy,
		DeletedAt:      deletedAt,
	}
}

func normalizeFarmerIdentityKey(name, community, nationalID string) string {
	return strings.ToLower(strings.TrimSpace(name)) + "|" +
		strings.ToLower(strings.TrimSpace(community)) + "|" +
		strings.ToLower(strings.TrimSpace(nationalID))
}

func shouldSkipSyncUpsert(existsByID, existsByIdentity bool, incomingID, conflictingID string) bool {
	return existsByIdentity && (!existsByID || conflictingID != incomingID)
}
