-- name: ListUsers :many
SELECT app_user_id, username, can_create, can_update, can_delete, is_meta
FROM app_user
ORDER BY app_user_id;

-- name: CreateUser :one
INSERT INTO app_user (username, password_hash, can_create, can_update, can_delete, is_meta)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING app_user_id, username, can_create, can_update, can_delete, is_meta;

-- name: UpdateUser :one
UPDATE app_user
SET
    username   = COALESCE(sqlc.narg(username),   username),
    can_create = COALESCE(sqlc.narg(can_create), can_create),
    can_update = COALESCE(sqlc.narg(can_update), can_update),
    can_delete = COALESCE(sqlc.narg(can_delete), can_delete),
    is_meta    = COALESCE(sqlc.narg(is_meta),    is_meta)
WHERE app_user_id = sqlc.arg(app_user_id)
RETURNING app_user_id, username, can_create, can_update, can_delete, is_meta;

-- name: DeleteUser :exec
DELETE FROM app_user
WHERE app_user_id = $1;
