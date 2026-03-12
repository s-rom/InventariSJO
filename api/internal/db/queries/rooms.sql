-- name: ListRoomsByCenter :many
SELECT *
FROM room
WHERE center_id = $1
ORDER BY room_id;

-- name: CreateRoom :one
INSERT INTO room (center_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateRoom :one
UPDATE room
SET name = $1
WHERE room_id = $2
RETURNING *;

-- name: DeleteRoom :exec
DELETE FROM room
WHERE room_id = $1;
