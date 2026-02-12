ALTER TABLE farmers
ADD COLUMN IF NOT EXISTS deleted_at BIGINT NULL;

CREATE INDEX IF NOT EXISTS farmers_deleted_at_idx
ON farmers (deleted_at);