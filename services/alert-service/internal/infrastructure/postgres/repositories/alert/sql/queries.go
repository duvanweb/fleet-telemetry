package sql

const InsertAlert = `
INSERT INTO alerts (id, vehicle_id, type, status, started_at, resolved_at, metadata, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, vehicle_id, type, status, started_at, resolved_at, metadata, created_at, updated_at`

const SelectAllAlerts = `
SELECT id, vehicle_id, type, status, started_at, resolved_at, metadata, created_at, updated_at
FROM alerts
ORDER BY created_at DESC`

const SelectOpenByVehicle = `
SELECT id, vehicle_id, type, status, started_at, resolved_at, metadata, created_at, updated_at
FROM alerts
WHERE vehicle_id = $1 AND status = 'OPEN'
LIMIT 1`

const UpdateMarkResolved = `
UPDATE alerts SET status = 'RESOLVED', resolved_at = NOW(), updated_at = NOW()
WHERE id = $1`
