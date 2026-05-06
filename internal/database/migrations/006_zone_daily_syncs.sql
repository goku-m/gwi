CREATE TABLE IF NOT EXISTS zone_daily_syncs (
    zone_name TEXT NOT NULL,
    sync_date DATE NOT NULL DEFAULT CURRENT_DATE,
    sync_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (zone_name, sync_date)
);

