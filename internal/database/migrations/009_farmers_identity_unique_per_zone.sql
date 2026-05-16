-- De-duplicate active farmers by zone + (name, community, national_id) before adding uniqueness guard.
-- Keep the oldest row (created_at, then id) and soft-delete the rest.
WITH ranked AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY
                LOWER(BTRIM(zone_name)),
                LOWER(BTRIM(name)),
                LOWER(BTRIM(community)),
                LOWER(BTRIM(COALESCE(national_id, '')))
            ORDER BY created_at ASC, id ASC
        ) AS rn
    FROM farmers
    WHERE deleted_at IS NULL
)
UPDATE farmers f
SET
    deleted_at = (EXTRACT(EPOCH FROM NOW()) * 1000)::bigint,
    updated_at = NOW()
FROM ranked r
WHERE f.id = r.id
  AND r.rn > 1
  AND f.deleted_at IS NULL;

-- Enforce in DB: one active farmer identity per zone (case/space-insensitive).
CREATE UNIQUE INDEX IF NOT EXISTS idx_farmers_zone_identity_unique_active
ON farmers (
    LOWER(BTRIM(zone_name)),
    LOWER(BTRIM(name)),
    LOWER(BTRIM(community)),
    LOWER(BTRIM(COALESCE(national_id, '')))
)
WHERE deleted_at IS NULL;
