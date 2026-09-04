CREATE TABLE IF NOT EXISTS telemetry_points (
    id                VARCHAR(26)      PRIMARY KEY,
    vehicle_id        VARCHAR(26)      NOT NULL,
    latitude          DOUBLE PRECISION NOT NULL,
    longitude         DOUBLE PRECISION NOT NULL,
    device_timestamp  TIMESTAMPTZ      NOT NULL,
    received_at       TIMESTAMPTZ      NOT NULL DEFAULT now(),
    deduplication_key VARCHAR(64)      NOT NULL UNIQUE,
    created_at        TIMESTAMPTZ      NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_telemetry_vid_dts ON telemetry_points(vehicle_id, device_timestamp);
CREATE INDEX IF NOT EXISTS idx_telemetry_vid_rat ON telemetry_points(vehicle_id, received_at);
