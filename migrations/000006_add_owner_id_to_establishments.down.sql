DROP INDEX IF EXISTS idx_establishments_owner_id;
ALTER TABLE establishments DROP CONSTRAINT IF EXISTS fk_establishments_owner;
ALTER TABLE establishments DROP COLUMN IF EXISTS owner_id;
