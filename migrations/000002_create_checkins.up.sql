-- Fase 3: tabela de checkins

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE checkins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    establishment_id UUID NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,
    checked_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_checkins_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_checkins_establishment FOREIGN KEY (establishment_id) REFERENCES establishments(id) ON DELETE CASCADE
);

-- coluna geography gerada para index espacial
ALTER TABLE checkins ADD COLUMN location geography(Point,4326) GENERATED ALWAYS AS (ST_SetSRID(ST_MakePoint(lng, lat), 4326)::geography) STORED;

CREATE INDEX idx_checkins_location ON checkins USING GIST (location);

-- índice para buscas por usuário/estabelecimento por dia
CREATE INDEX idx_checkins_user_estab_checked_at ON checkins (user_id, establishment_id, checked_at);
