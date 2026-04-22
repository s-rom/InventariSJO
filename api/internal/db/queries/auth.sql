-- name: GetUserByUsername :one
SELECT *
FROM app_user
WHERE username = $1;

-- name: GetUserByGoogleSub :one
SELECT *
FROM app_user
WHERE google_sub = $1;

-- name: GetUserByEmail :one
SELECT *
FROM app_user
WHERE email = $1;

-- name: CreateGoogleUser :one
INSERT INTO app_user (username, email, google_sub, role_id)
VALUES ($1, $2, $3, 'readonly')
RETURNING *;

-- name: SetGoogleSub :one
UPDATE app_user
SET google_sub = $1
WHERE email = $2
RETURNING *;
