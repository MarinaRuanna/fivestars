ALTER TABLE establishments
DROP COLUMN IF EXISTS claim_code_expires_at,
DROP COLUMN IF EXISTS claim_code_hash;
