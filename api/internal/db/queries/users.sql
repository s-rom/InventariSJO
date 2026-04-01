-- name: ListUsers :many
SELECT app_user_id, username, role_id
FROM app_user
ORDER BY app_user_id;

-- name: CreateUser :one
INSERT INTO app_user (username, password_hash, role_id)
VALUES (@username, @password_hash, @role_id)
RETURNING app_user_id, username, role_id;

-- name: UpdateUser :one
UPDATE app_user
SET
    username = COALESCE(sqlc.narg(username), username),
    role_id  = COALESCE(sqlc.narg(role_id),  role_id)
WHERE app_user_id = sqlc.arg(app_user_id)
RETURNING app_user_id, username, role_id;

-- name: UpdateUserPassword :exec
UPDATE app_user
SET password_hash = $2
WHERE app_user_id = $1;

-- name: DeleteUser :exec
DELETE FROM app_user
WHERE app_user_id = $1;
