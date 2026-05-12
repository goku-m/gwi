package farmer

import (
	"time"

	"github.com/goku-m/gwi/internal/model"
)

// --------------------------------------------------
// Core Farmer Model
// --------------------------------------------------

type Farmer struct {
	model.Base

	ZoneName string `json:"zoneName" db:"zone_name"`

	Name       string `json:"name" db:"name"`
	NationalID string `json:"nationalId" db:"national_id"`
	Community  string `json:"community" db:"community"`

	Prefinance float64 `json:"prefinance" db:"prefinance"`
	Balance    float64 `json:"balance" db:"balance"`

	TotalKgBrought float64 `json:"totalKgBrought" db:"total_kg_brought"`
	TotalAmount    float64 `json:"totalAmount" db:"total_amount"`
	CreatedBy      string  `json:"createdBy" db:"created_by"`
	UpdatedBy      string  `json:"updatedBy" db:"updated_by"`

	// Soft-delete tombstone timestamp (epoch ms).
	DeletedAt *int64 `json:"deletedAt,omitempty" db:"deleted_at"`
}

// --------------------------------------------------
// Expanded / Joined Variant (future-proof)
// --------------------------------------------------

type PopulatedFarmer struct {
	Farmer
	// Add related data later:
	// Deliveries []Delivery `json:"deliveries"`
	// Zone       *Zone      `json:"zone"`
}

type EditStatus struct {
	ShouldEdit bool `json:"shouldEdit" db:"should_edit"`
	// Add related data later:
	// Deliveries []Delivery `json:"deliveries"`
	// Zone       *Zone      `json:"zone"`
}

// --------------------------------------------------
// Aggregates / Stats
// --------------------------------------------------

type FarmerStats struct {
	ZoneName string `json:"zoneName" db:"zone_name"`

	TotalFarmers     int     `json:"totalFarmers" db:"total_farmers"`
	TotalCommunities int     `json:"totalCommunities" db:"total_communities"`
	DailySyncs       int     `json:"dailySyncs" db:"daily_syncs"`
	TotalKgBrought   float64 `json:"totalKgBrought" db:"total_kg_brought"`
	TotalAmount      float64 `json:"totalAmount" db:"total_amount"`
	TotalPrefinance  float64 `json:"totalPrefinance" db:"total_prefinance"`
	TotalBalance     float64 `json:"totalBalance" db:"total_balance"` // optional but useful
}

type CommunityFarmerStats struct {
	ZoneName      string `json:"zoneName" db:"zone_name"`
	CommunityName string `json:"communityName" db:"community_name"`

	TotalFarmers    int     `json:"totalFarmers" db:"total_farmers"`
	DailySyncs      int     `json:"dailySyncs" db:"daily_syncs"`
	TotalKgBrought  float64 `json:"totalKgBrought" db:"total_kg_brought"`
	TotalAmount     float64 `json:"totalAmount" db:"total_amount"`
	TotalPrefinance float64 `json:"totalPrefinance" db:"total_prefinance"`
	TotalBalance    float64 `json:"totalBalance" db:"total_balance"`
}

type ZoneCommunitiesResponse struct {
	ZoneName    string   `json:"zoneName"`
	Communities []string `json:"communities"`
}

type NewFarmersStats struct {
	NewFarmers int    `json:"newFarmers" db:"new_farmers"`
	SinceDate  string `json:"sinceDate" db:"since_date"`
}

// --------------------------------------------------
// Helper methods
// --------------------------------------------------

func (f *Farmer) HasOutstandingBalance() bool {
	return f.Balance > 0
}

func (f *Farmer) LastUpdated() time.Time {
	return f.UpdatedAt
}
