package sql

const InsertOutboxEvent = `
INSERT INTO outbox_events (id, event_type, payload, status, created_at)
VALUES ($1, $2, $3, $4, $5)`

const SelectPending = `
SELECT id, event_type, payload, status, created_at, published_at
FROM outbox_events
WHERE status = 'pending'
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED`

const UpdateMarkPublished = `
UPDATE outbox_events
SET status = 'published', published_at = NOW()
WHERE id = ANY($1)`
