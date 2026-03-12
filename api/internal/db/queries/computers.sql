-- name: ListComputers :many
-- Returns all computers with computer_type derived from which subtable has a row.
SELECT
    c.*,
    CASE WHEN d.computer_id IS NOT NULL THEN 'desktop' ELSE 'laptop' END AS computer_type
FROM computer c
LEFT JOIN desktop d ON d.computer_id = c.computer_id
ORDER BY c.computer_id;

-- name: GetComputerBase :one
SELECT * FROM computer WHERE computer_id = $1;

-- name: CreateComputer :one
INSERT INTO computer (hostname, room_id, observations, created_by_app_user_id)
VALUES (@hostname, sqlc.narg(room_id), sqlc.narg(observations), @created_by_app_user_id)
RETURNING *;

-- name: UpdateComputerBase :one
UPDATE computer
SET
    hostname     = COALESCE(sqlc.narg(hostname),     hostname),
    room_id      = COALESCE(sqlc.narg(room_id),      room_id),
    observations = COALESCE(sqlc.narg(observations),  observations),
    updated_at   = now()
WHERE computer_id = sqlc.arg(computer_id)
RETURNING *;

-- name: DeleteComputer :exec
DELETE FROM computer WHERE computer_id = $1;

