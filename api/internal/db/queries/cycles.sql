-- name: ListCycles :many
SELECT * FROM cycle ORDER BY name;

-- name: CreateCycle :one
INSERT INTO cycle (name) VALUES ($1) RETURNING *;

-- name: UpdateCycle :one
UPDATE cycle SET name = $2 WHERE cycle_id = $1 RETURNING *;

-- name: DeleteCycle :exec
DELETE FROM cycle WHERE cycle_id = $1;
