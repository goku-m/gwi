package repository

import (
	"testing"
	"time"
)

func TestNormalizeCreatedBy(t *testing.T) {
	if got := normalizeCreatedBy("  Alice  "); got != "Alice" {
		t.Fatalf("expected Alice, got %q", got)
	}
	if got := normalizeCreatedBy("   "); got != "UNKNOWN" {
		t.Fatalf("expected UNKNOWN, got %q", got)
	}
}

func TestNormalizeUpdatedBy(t *testing.T) {
	if got := normalizeUpdatedBy("  Bob  "); got != "Bob" {
		t.Fatalf("expected Bob, got %q", got)
	}
	if got := normalizeUpdatedBy(""); got != "NIL" {
		t.Fatalf("expected NIL, got %q", got)
	}
}

func TestTimeMsRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	ms := timeToMs(now)
	back := msToTime(ms)
	if !back.Equal(now) {
		t.Fatalf("expected %v, got %v", now, back)
	}
}

func TestFarmerSyncFromRowMapsAddedUpdatedBy(t *testing.T) {
	now := time.Now().UTC()
	row := SyncFarmerRow{
		ID:            "id-1",
		ZoneName:      "ZONE",
		Name:          "Farmer",
		NationalID:    "N1",
		Community:     "C1",
		Prefinance:    10,
		Balance:       5,
		TotalKg:       2,
		TotalAmount:   20,
		CreatedAt:     now,
		UpdatedAt:     now,
		CreatedBy:     "Creator",
		UpdatedBy:     "Updater",
		DeletedAtNull: true,
	}

	record := farmerSyncFromRow(row)
	if record.AddedBy != "Creator" {
		t.Fatalf("expected AddedBy Creator, got %q", record.AddedBy)
	}
	if record.UpdatedBy != "Updater" {
		t.Fatalf("expected UpdatedBy Updater, got %q", record.UpdatedBy)
	}
	if record.DeletedAt != nil {
		t.Fatalf("expected nil DeletedAt, got %v", *record.DeletedAt)
	}
}
