CREATE TABLE IF NOT EXISTS alerts (
    id          VARCHAR(26)  PRIMARY KEY,
    vehicle_id  VARCHAR(26)  NOT NULL,
    type        VARCHAR(50)  NOT NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'OPEN',
    started_at  TIMESTAMPTZ  NOT NULL,
    resolved_at TIMESTAMPTZ,
    metadata    JSONB,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_alerts_vehicle_status ON alerts(vehicle_id, status);
CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status);
