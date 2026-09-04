CREATE TABLE IF NOT EXISTS vehicles (
    id          VARCHAR(26)  PRIMARY KEY,
    external_id VARCHAR(255) NOT NULL UNIQUE,
    plate       VARCHAR(50)  NOT NULL UNIQUE,
    name        VARCHAR(255) NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
