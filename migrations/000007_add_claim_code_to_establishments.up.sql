ALTER TABLE establishments
ADD COLUMN claim_code_hash TEXT NULL,
ADD COLUMN claim_code_expires_at TIMESTAMPTZ NULL;
