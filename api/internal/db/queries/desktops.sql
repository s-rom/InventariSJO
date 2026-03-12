-- name: ListDesktops :many
SELECT
    c.computer_id, c.hostname, c.room_id, c.observations,
    c.created_by_app_user_id, c.created_at, c.updated_at,
    d.desktop_model_id, d.cpu_id, d.ram_gb, d.ram_type,
    d.storage_gb, d.storage_type, d.os_id,
    d.equipment_user_id, d.has_wifi_card, d.mac_address
FROM computer c
INNER JOIN desktop d ON d.computer_id = c.computer_id
ORDER BY c.computer_id;

-- name: GetDesktop :one
SELECT
    c.computer_id, c.hostname, c.room_id, c.observations,
    c.created_by_app_user_id, c.created_at, c.updated_at,
    d.desktop_model_id, d.cpu_id, d.ram_gb, d.ram_type,
    d.storage_gb, d.storage_type, d.os_id,
    d.equipment_user_id, d.has_wifi_card, d.mac_address
FROM computer c
INNER JOIN desktop d ON d.computer_id = c.computer_id
WHERE c.computer_id = $1;

-- name: CreateDesktop :one
INSERT INTO desktop (
    computer_id, desktop_model_id, cpu_id,
    ram_gb, ram_type, storage_gb, storage_type,
    os_id, equipment_user_id, has_wifi_card, mac_address
) VALUES (
    @computer_id, sqlc.narg(desktop_model_id), sqlc.narg(cpu_id),
    sqlc.narg(ram_gb), sqlc.narg(ram_type), sqlc.narg(storage_gb), sqlc.narg(storage_type),
    sqlc.narg(os_id), sqlc.narg(equipment_user_id), @has_wifi_card, sqlc.narg(mac_address)
) RETURNING *;

-- name: UpdateDesktop :one
UPDATE desktop
SET
    desktop_model_id  = COALESCE(sqlc.narg(desktop_model_id),  desktop_model_id),
    cpu_id            = COALESCE(sqlc.narg(cpu_id),            cpu_id),
    ram_gb            = COALESCE(sqlc.narg(ram_gb),            ram_gb),
    ram_type          = COALESCE(sqlc.narg(ram_type),          ram_type),
    storage_gb        = COALESCE(sqlc.narg(storage_gb),        storage_gb),
    storage_type      = COALESCE(sqlc.narg(storage_type),      storage_type),
    os_id             = COALESCE(sqlc.narg(os_id),             os_id),
    equipment_user_id = COALESCE(sqlc.narg(equipment_user_id), equipment_user_id),
    has_wifi_card     = COALESCE(sqlc.narg(has_wifi_card),     has_wifi_card),
    mac_address       = COALESCE(sqlc.narg(mac_address),       mac_address)
WHERE computer_id = sqlc.arg(computer_id)
RETURNING *;
