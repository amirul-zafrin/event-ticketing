-- name: CreatePrice :one
INSERT INTO prices (
  class, price, event
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetPrice :one
SELECT * FROM prices
WHERE id = $1;

-- name: ListAllPrices :many
SELECT * FROM prices
LIMIT $1
OFFSET $2;

-- name: UpdatePrice :one
UPDATE prices
SET price = $2,
    updated_at = NOW()
WHERE id =  $1
RETURNING *;

-- name: SoftDeletePrice :exec
UPDATE prices
SET deleted_at = NOW()
WHERE id = $1;

-- name: PermaDeletePrice :exec
DELETE FROM prices
WHERE id = $1;