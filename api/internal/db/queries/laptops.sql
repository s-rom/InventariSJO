-- name: ListLaptops :many
SELECT
    c.computer_id, c.hostname, c.room_id, c.observations,
    c.created_by_app_user_id, c.created_at, c.updated_at,
    l.laptop_model_id, l.ram_gb, l.ram_type,
    l.storage_gb, l.storage_type, l.mac_address,
    l.os_id, l.equipment_user_id
FROM computer c
INNER JOIN laptop l ON l.computer_id = c.computer_id
ORDER BY c.computer_id;

-- name: GetLaptop :one
SELECT
    c.computer_id, c.hostname, c.room_id, c.observations,
    c.created_by_app_user_id, c.created_at, c.updated_at,
    l.laptop_model_id, l.ram_gb, l.ram_type,
    l.storage_gb, l.storage_type, l.mac_address,
    l.os_id, l.equipment_user_id
FROM computer c
INNER JOIN laptop l ON l.computer_id = c.computer_id
WHERE c.computer_id = $1;

-- name: CreateLaptop :one
INSERT INTO laptop (
    computer_id, laptop_model_id,
    ram_gb, ram_type, storage_gb, storage_type,
    mac_address, os_id, equipment_user_id
) VALUES (
    @computer_id, @laptop_model_id,
    sqlc.narg(ram_gb), sqlc.narg(ram_type), sqlc.narg(storage_gb), sqlc.narg(storage_type),
    sqlc.narg(mac_address), sqlc.narg(os_id), sqlc.narg(equipment_user_id)
) RETURNING *;

-- name: UpdateLaptop :one
UPDATE laptop
SET
    laptop_model_id   = COALESCE(sqlc.narg(laptop_model_id),   laptop_model_id),
    ram_gb            = COALESCE(sqlc.narg(ram_gb),            ram_gb),
    ram_type          = COALESCE(sqlc.narg(ram_type),          ram_type),
    storage_gb        = COALESCE(sqlc.narg(storage_gb),        storage_gb),
    storage_type      = COALESCE(sqlc.narg(storage_type),      storage_type),
    mac_address       = COALESCE(sqlc.narg(mac_address),       mac_address),
    os_id             = COALESCE(sqlc.narg(os_id),             os_id),
    equipment_user_id = COALESCE(sqlc.narg(equipment_user_id), equipment_user_id)
WHERE computer_id = sqlc.arg(computer_id)
RETURNING *;
