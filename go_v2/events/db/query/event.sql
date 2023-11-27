-- name: CreateEvent :one
INSERT INTO events (
  name, date, location, capacity
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateEventSeat :one
UPDATE events
SET seats = $2
WHERE id = $1
RETURNING *;

-- name: GetEvent :one
SELECT * FROM events
WHERE id = $1 LIMIT 1;

-- name: ListAllEvents :many
SELECT * FROM events
ORDER by id
LIMIT $1
OFFSET $2;

-- name: UpdateEvent :one
UPDATE events
SET name = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteEvent :exec
UPDATE events
SET deleted_at = NOW()
WHERE id = $1;

-- name: PermaDeleteEvent :exec
DELETE FROM events
WHERE id = $1;