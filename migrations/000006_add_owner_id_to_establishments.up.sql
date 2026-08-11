ALTER TABLE establishments
ADD COLUMN owner_id UUID NULL,
ADD CONSTRAINT fk_establishments_owner
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL;

CREATE INDEX idx_establishments_owner_id ON establishments(owner_id);
