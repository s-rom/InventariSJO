-- name: ListEquipmentUsers :many
SELECT *
FROM equipment_user
ORDER BY equipment_user_id;

-- name: CreateEquipmentUser :one
INSERT INTO equipment_user (name)
VALUES ($1)
RETURNING *;

-- name: UpdateEquipmentUser :one
UPDATE equipment_user
SET name = $1
WHERE equipment_user_id = $2
RETURNING *;

-- name: DeleteEquipmentUser :exec
DELETE FROM equipment_user
WHERE equipment_user_id = $1;
