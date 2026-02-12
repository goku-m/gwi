
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =========================================================
-- Farmers table (zone-specific, single shared table)
-- =========================================================
CREATE TABLE edit_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Permanent zone name (hardcoded in zone-specific apps)
    should_edit BOOLEAN  NOT NULL DEFAULT TRUE,

   

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);