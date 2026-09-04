package sql

const InsertVehicle = `
INSERT INTO vehicles (id, external_id, plate, name, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, external_id, plate, name, created_at, updated_at, deleted_at`

const ExistsByExternalID = `
SELECT EXISTS(SELECT 1 FROM vehicles WHERE external_id = $1 AND deleted_at IS NULL)`

const ExistsByPlate = `
SELECT EXISTS(SELECT 1 FROM vehicles WHERE plate = $1 AND deleted_at IS NULL)`

const SelectAll = `
SELECT id, external_id, plate, name, created_at, updated_at, deleted_at
FROM vehicles
WHERE deleted_at IS NULL
ORDER BY created_at DESC`

const SelectByID = `
SELECT id, external_id, plate, name, created_at, updated_at, deleted_at
FROM vehicles
WHERE id = $1`

const SoftDelete = `
UPDATE vehicles SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1 AND deleted_at IS NULL`

const UpdateVehicle = `
UPDATE vehicles SET plate = $1, name = $2, updated_at = NOW()
WHERE id = $3 AND deleted_at IS NULL
RETURNING id, external_id, plate, name, created_at, updated_at, deleted_at`
