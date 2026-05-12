CREATE TABLE IF NOT EXISTS farmer_deletions (
    farmer_id UUID NOT NULL,
    zone_name TEXT NOT NULL,
    deleted_at BIGINT NOT NULL,
    PRIMARY KEY (farmer_id, zone_name)
);

CREATE INDEX IF NOT EXISTS idx_farmer_deletions_zone_deleted_at
ON farmer_deletions (zone_name, deleted_at);

CREATE OR REPLACE FUNCTION record_farmer_deletion()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO farmer_deletions (farmer_id, zone_name, deleted_at)
    VALUES (
        OLD.id,
        OLD.zone_name,
        COALESCE(OLD.deleted_at, (EXTRACT(EPOCH FROM now()) * 1000)::bigint)
    )
    ON CONFLICT (farmer_id, zone_name) DO UPDATE
    SET deleted_at = EXCLUDED.deleted_at;

    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_record_farmer_deletion ON farmers;
CREATE TRIGGER trg_record_farmer_deletion
AFTER DELETE ON farmers
FOR EACH ROW
EXECUTE FUNCTION record_farmer_deletion();

