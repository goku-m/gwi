-- Required extension for UUIDs
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =========================================================
-- Farmers table (zone-specific, single shared table)
-- =========================================================
CREATE TABLE farmers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Permanent zone name (hardcoded in zone-specific apps)
    zone_name TEXT NOT NULL,

    name TEXT NOT NULL,
    national_id TEXT  ,
    community TEXT NOT NULL,

    prefinance NUMERIC(12,2) NOT NULL DEFAULT 0,
    balance NUMERIC(12,2) NOT NULL DEFAULT 0,

    total_kg_brought NUMERIC(12,2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(12,2) NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =========================================================
-- Enforce allowed permanent zones (edit list as needed)
-- =========================================================
ALTER TABLE farmers
ADD CONSTRAINT farmers_zone_name_check
CHECK (
    zone_name IN (
        'WA',
        'YENDI',
        'GARU',
        'NAPKADURI',
        'LANGBINSI'
    )
);

-- =========================================================
-- Indexes for fast zone-based CRUD and reporting
-- =========================================================
CREATE INDEX idx_farmers_zone_name ON farmers(zone_name);
CREATE INDEX idx_farmers_zone_community ON farmers(zone_name, community);
CREATE INDEX idx_farmers_national_id ON farmers(national_id);

-- =========================================================
-- Auto-update updated_at timestamp
-- (assumes trigger_set_updated_at() already exists)
-- =========================================================
CREATE TRIGGER set_updated_at_farmers
    BEFORE UPDATE ON farmers
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_updated_at();



