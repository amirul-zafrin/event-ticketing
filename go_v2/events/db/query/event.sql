-- name: CreateEvent :one
INSERT INTO events (
  name, date, location, capacity
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;
