-- name: ListCenters :many
SELECT *
FROM center
ORDER BY center_id;

-- name: CreateCenter :one
INSERT INTO center (name)
VALUES ($1)
RETURNING *;

-- name: UpdateCenter :one
UPDATE center
SET name = $1
WHERE center_id = $2
RETURNING *;

-- name: DeleteCenter :exec
DELETE FROM center
WHERE center_id = $1;
