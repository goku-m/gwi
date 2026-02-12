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
	conditions := []string{"t.zone_name = @zone_name"}

	if query != nil {
		if query.Search != nil && strings.TrimSpace(*query.Search) != "" {
			conditions = append(conditions, "(t.name ILIKE @search OR t.national_id ILIKE @search)")
			args["search"] = "%" + strings.TrimSpace(*query.Search) + "%"
		}

		if query.Community != nil && strings.TrimSpace(*query.Community) != "" {
			conditions = append(conditions, "t.community ILIKE @community")
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

func (r *FarmerRepository) GetZoneStats(ctx context.Context, zoneName string) (*farmer.FarmerStats, error) {
	stmt := `
		SELECT
			@zone_name::text AS zone_name,
			COUNT(*) AS total_farmers,
			COALESCE(SUM(total_kg_brought), 0)::float8 AS total_kg_brought,
			COALESCE(SUM(total_amount), 0)::float8 AS total_amount,
			COALESCE(SUM(prefinance), 0)::float8 AS total_prefinance,
			COALESCE(SUM(balance), 0)::float8 AS total_balance
		FROM
			farmers
		WHERE
			zone_name = @zone_name
	`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"zone_name": zoneName,
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
				&row.CreatedAt, &row.UpdatedAt, &row.DeletedAtNull, &row.DeletedAtMs,
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
		SELECT id::text
		FROM farmers
		WHERE zone_name = @zone_name
		  AND deleted_at IS NOT NULL
		  AND deleted_at > @last_ms
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
				&row.CreatedAt, &row.UpdatedAt, &row.DeletedAtNull, &row.DeletedAtMs,
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
				&row.CreatedAt, &row.UpdatedAt, &row.DeletedAtNull, &row.DeletedAtMs,
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
			created_at, updated_at, deleted_at
		) VALUES (
			@id, @zone_name, @name, @national_id, @community,
			@prefinance, @balance, @total_kg_brought, @total_amount,
			now(), now(), NULL
		)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			national_id = EXCLUDED.national_id,
			community = EXCLUDED.community,
			prefinance = EXCLUDED.prefinance,
			balance = EXCLUDED.balance,
			total_kg_brought = EXCLUDED.total_kg_brought,
			total_amount = EXCLUDED.total_amount,
			updated_at = now(),
			deleted_at = NULL
		WHERE farmers.zone_name = EXCLUDED.zone_name
	`

	deleteSQL := `
		UPDATE farmers
		SET deleted_at = (EXTRACT(EPOCH FROM now()) * 1000)::bigint, updated_at = now()
		WHERE id = @id AND zone_name = @zone_name
	`

	b := &pgx.Batch{}

	for _, f := range upserts {
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
		})
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
		DeletedAt:      deletedAt,
	}
}
