-- name: ListDesktops :many
SELECT
    c.computer_id, c.hostname, c.room_id, c.observations,
    c.created_by_app_user_id, c.created_at, c.updated_at,
    d.desktop_model_id,
    COALESCE(d.cpu_id,       CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.cpu_id       END) AS cpu_id,
    COALESCE(d.ram_gb,       CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_ram_gb       END) AS ram_gb,
    COALESCE(d.ram_type,     CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_ram_type     END) AS ram_type,
    COALESCE(d.storage_gb,   CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_storage_gb   END) AS storage_gb,
    COALESCE(d.storage_type, CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_storage_type END) AS storage_type,
    COALESCE(d.os_id,        CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_os_id        END) AS os_id,
    d.equipment_user_id, d.has_wifi_card, d.mac_address, d.network_connection
FROM computer c
INNER JOIN desktop d ON d.computer_id = c.computer_id
LEFT JOIN desktop_model dm ON dm.desktop_model_id = d.desktop_model_id
ORDER BY c.computer_id;

-- name: GetDesktop :one
SELECT
    c.computer_id, c.hostname, c.room_id, c.observations,
    c.created_by_app_user_id, c.created_at, c.updated_at,
    d.desktop_model_id,
    COALESCE(d.cpu_id,       CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.cpu_id       END) AS cpu_id,
    COALESCE(d.ram_gb,       CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_ram_gb       END) AS ram_gb,
    COALESCE(d.ram_type,     CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_ram_type     END) AS ram_type,
    COALESCE(d.storage_gb,   CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_storage_gb   END) AS storage_gb,
    COALESCE(d.storage_type, CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_storage_type END) AS storage_type,
    COALESCE(d.os_id,        CASE WHEN dm.desktop_model_id IS NOT NULL THEN dm.base_os_id        END) AS os_id,
    d.equipment_user_id, d.has_wifi_card, d.mac_address, d.network_connection
FROM computer c
INNER JOIN desktop d ON d.computer_id = c.computer_id
LEFT JOIN desktop_model dm ON dm.desktop_model_id = d.desktop_model_id
WHERE c.computer_id = $1;

-- name: CreateDesktop :one
INSERT INTO desktop (
    computer_id, desktop_model_id, cpu_id,
    ram_gb, ram_type, storage_gb, storage_type,
    os_id, equipment_user_id, has_wifi_card, mac_address, network_connection
) VALUES (
    @computer_id, sqlc.narg(desktop_model_id), sqlc.narg(cpu_id),
    sqlc.narg(ram_gb), sqlc.narg(ram_type), sqlc.narg(storage_gb), sqlc.narg(storage_type),
    sqlc.narg(os_id), sqlc.narg(equipment_user_id), @has_wifi_card, sqlc.narg(mac_address), sqlc.narg(network_connection)
) RETURNING *;

-- name: UpdateDesktop :one
UPDATE desktop
SET
    desktop_model_id   = COALESCE(sqlc.narg(desktop_model_id),   desktop_model_id),
    cpu_id             = COALESCE(sqlc.narg(cpu_id),             cpu_id),
    ram_gb             = COALESCE(sqlc.narg(ram_gb),             ram_gb),
    ram_type           = COALESCE(sqlc.narg(ram_type),           ram_type),
    storage_gb         = COALESCE(sqlc.narg(storage_gb),         storage_gb),
    storage_type       = COALESCE(sqlc.narg(storage_type),       storage_type),
    os_id              = COALESCE(sqlc.narg(os_id),              os_id),
    equipment_user_id  = COALESCE(sqlc.narg(equipment_user_id),  equipment_user_id),
    has_wifi_card      = COALESCE(sqlc.narg(has_wifi_card),      has_wifi_card),
    mac_address        = COALESCE(sqlc.narg(mac_address),        mac_address),
    network_connection = COALESCE(sqlc.narg(network_connection), network_connection)
WHERE computer_id = sqlc.arg(computer_id)
RETURNING *;
