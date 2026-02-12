package farmer

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ------------------------------------------------------------
// Create
// ------------------------------------------------------------

type CreateFarmerPayload struct {
	Name       string `json:"name" validate:"required,min=1,max=255"`
	NationalID string `json:"nationalId" validate:"required,min=3,max=50"`
	Community  string `json:"community" validate:"required,min=1,max=255"`

	// Optional on create (default 0 in DB); allow if you want to set from app.
	Prefinance      *float64 `json:"prefinance" validate:"omitempty,gte=0"`
	Balance         *float64 `json:"balance" validate:"omitempty,gte=0"`
	TotalKgBrought  *float64 `json:"totalKgBrought" validate:"omitempty,gte=0"`
	TotalAmount     *float64 `json:"totalAmount" validate:"omitempty,gte=0"`
}

func (p *CreateFarmerPayload) Validate() error {
	validate := validator.New()

	// Normalize common fields (optional but helpful)
	p.Name = strings.TrimSpace(p.Name)
	p.NationalID = strings.TrimSpace(p.NationalID)
	p.Community = strings.TrimSpace(p.Community)

	return validate.Struct(p)
}

// ------------------------------------------------------------
// Update
// ------------------------------------------------------------

type UpdateFarmerPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`

	Name       *string `json:"name" validate:"omitempty,min=1,max=255"`
	NationalID *string `json:"nationalId" validate:"omitempty,min=3,max=50"`
	Community  *string `json:"community" validate:"omitempty,min=1,max=255"`

	Prefinance     *float64 `json:"prefinance" validate:"omitempty,gte=0"`
	Balance        *float64 `json:"balance" validate:"omitempty,gte=0"`
	TotalKgBrought *float64 `json:"totalKgBrought" validate:"omitempty,gte=0"`
	TotalAmount    *float64 `json:"totalAmount" validate:"omitempty,gte=0"`
}

func (p *UpdateFarmerPayload) Validate() error {
	validate := validator.New()

	if p.Name != nil {
		s := strings.TrimSpace(*p.Name)
		p.Name = &s
	}
	if p.NationalID != nil {
		s := strings.TrimSpace(*p.NationalID)
		p.NationalID = &s
	}
	if p.Community != nil {
		s := strings.TrimSpace(*p.Community)
		p.Community = &s
	}

	return validate.Struct(p)
}

// ------------------------------------------------------------
// List Query (per zone)
// ------------------------------------------------------------

type GetFarmersQuery struct {
	Page  *int    `query:"page" validate:"omitempty,min=1"`
	Limit *int    `query:"limit" validate:"omitempty,min=1,max=100"`

	// Sorting allowed fields based on your farmers table
	Sort  *string `query:"sort" validate:"omitempty,oneof=created_at updated_at name community national_id balance total_kg_brought total_amount"`
	Order *string `query:"order" validate:"omitempty,oneof=asc desc"`

	Search    *string `query:"search" validate:"omitempty,min=1"`    // match name/national_id
	Community *string `query:"community" validate:"omitempty,min=1"` // filter by community
	HasDebt   *bool   `query:"hasDebt"`                               // true => balance > 0
}

type GetEditQuery struct {
	ShouldEdit   *bool   `query:"shouldEdit"`                               // true => balance > 0
}

func (q *GetFarmersQuery) Validate() error {
	validate := validator.New()

	if err := validate.Struct(q); err != nil {
		return err
	}

	// Defaults
	if q.Page == nil {
		v := 1
		q.Page = &v
	}
	if q.Limit == nil {
		v := 20
		q.Limit = &v
	}
	if q.Sort == nil {
		v := "created_at"
		q.Sort = &v
	}
	if q.Order == nil {
		v := "desc"
		q.Order = &v
	}

	// Normalize filters
	if q.Search != nil {
		s := strings.TrimSpace(*q.Search)
		q.Search = &s
	}
	if q.Community != nil {
		s := strings.TrimSpace(*q.Community)
		q.Community = &s
	}

	return nil
}

// ------------------------------------------------------------
// By ID
// ------------------------------------------------------------

type GetFarmerByIDPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *GetFarmerByIDPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// ------------------------------------------------------------
// Delete
// ------------------------------------------------------------

type DeleteFarmerPayload struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (p *DeleteFarmerPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}






type PullChangesResponse struct {
	Changes   map[string]any `json:"changes"`
	Timestamp int64          `json:"timestamp"`
}

type PushChangesRequest struct {
	Changes      map[string]TableChanges[FarmerSyncRecord] `json:"changes"`
	LastPulledAt *int64                                    `json:"lastPulledAt,omitempty"`
	ChunkIndex   *int                                      `json:"chunkIndex,omitempty"`
	ChunkCount   *int                                      `json:"chunkCount,omitempty"`
}


type FarmerSyncRecord struct {
	ID string `json:"id"`

	ZoneName   string `json:"zone_name"`
	Name       string `json:"name"`
	NationalID string `json:"national_id"`
	Community  string `json:"community"`

	Prefinance     float64 `json:"prefinance"`
	Balance        float64 `json:"balance"`
	TotalKgBrought float64 `json:"total_kg_brought"`
	TotalAmount    float64 `json:"total_amount"`

	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	DeletedAt *int64 `json:"deleted_at,omitempty"`
}

type TableChanges[T any] struct {
	Created []T      `json:"created"`
	Updated []T      `json:"updated"`
	Deleted []string `json:"deleted"`
}

type PullResult struct {
	Created   []FarmerSyncRecord
	Updated   []FarmerSyncRecord
	Deleted   []string
	Timestamp int64
}
