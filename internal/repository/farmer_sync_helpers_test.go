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

func TestNormalizeFarmerIdentityKey(t *testing.T) {
	got := normalizeFarmerIdentityKey("  John Doe ", " North ", " AB-123 ")
	want := "john doe|north|ab-123"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestShouldSkipSyncUpsert(t *testing.T) {
	cases := []struct {
		name             string
		existsByID       bool
		existsByIdentity bool
		incomingID       string
		conflictingID    string
		wantSkip         bool
	}{
		{
			name:             "allow when identity does not exist",
			existsByID:       false,
			existsByIdentity: false,
			incomingID:       "new-id",
			conflictingID:    "",
			wantSkip:         false,
		},
		{
			name:             "allow update when same id owns identity",
			existsByID:       true,
			existsByIdentity: true,
			incomingID:       "same-id",
			conflictingID:    "same-id",
			wantSkip:         false,
		},
		{
			name:             "skip when another id already has identity",
			existsByID:       false,
			existsByIdentity: true,
			incomingID:       "new-id",
			conflictingID:    "existing-id",
			wantSkip:         true,
		},
		{
			name:             "skip when incoming id exists but identity belongs to different id",
			existsByID:       true,
			existsByIdentity: true,
			incomingID:       "id-a",
			conflictingID:    "id-b",
			wantSkip:         true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := shouldSkipSyncUpsert(tc.existsByID, tc.existsByIdentity, tc.incomingID, tc.conflictingID)
			if got != tc.wantSkip {
				t.Fatalf("expected skip=%v, got %v", tc.wantSkip, got)
			}
		})
	}
}
