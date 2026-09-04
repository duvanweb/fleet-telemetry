package sql

const InsertOutboxEventInTx = `
INSERT INTO outbox_events (id, event_type, payload, status, created_at)
VALUES ($1, $2, $3, $4, $5)`

const InsertTelemetry = `
INSERT INTO telemetry_points
    (id, vehicle_id, latitude, longitude, device_timestamp, received_at, deduplication_key)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, vehicle_id, latitude, longitude, device_timestamp, received_at, deduplication_key`

const SelectByVehicleID = `
SELECT id, vehicle_id, latitude, longitude, device_timestamp, received_at, deduplication_key
FROM telemetry_points
WHERE vehicle_id = $1
ORDER BY device_timestamp ASC`
