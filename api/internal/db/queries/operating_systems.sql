-- name: ListOS :many
SELECT *
FROM os
ORDER BY os_id;

-- name: CreateOS :one
INSERT INTO os (name)
VALUES ($1)
RETURNING *;

-- name: UpdateOS :one
UPDATE os
SET name = $1
WHERE os_id = $2
RETURNING *;

-- name: DeleteOS :exec
DELETE FROM os
WHERE os_id = $1;
