-- name: ListRoles :many
SELECT * FROM role ORDER BY role_id;

-- name: CreateRole :one
INSERT INTO role (role_id, description)
VALUES (sqlc.arg(role_id), sqlc.narg(description))
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM role WHERE role_id = $1;
