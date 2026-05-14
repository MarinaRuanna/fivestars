-- Fase 5: highlights por estabelecimento

CREATE TABLE highlights (
    establishment_id UUID NOT NULL,
    review_id UUID NOT NULL,
    created_by_user_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_highlights PRIMARY KEY (establishment_id, review_id),
    CONSTRAINT fk_highlights_establishment FOREIGN KEY (establishment_id) REFERENCES establishments(id) ON DELETE CASCADE,
    CONSTRAINT fk_highlights_review FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    CONSTRAINT fk_highlights_created_by_user FOREIGN KEY (created_by_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_highlights_establishment_created_at
ON highlights (establishment_id, created_at DESC);

CREATE INDEX idx_highlights_review_id
ON highlights (review_id);
