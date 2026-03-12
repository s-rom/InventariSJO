-- name: ListComputers :many
SELECT *
FROM computer
ORDER BY computer_id;

-- name: GetComputer :one
SELECT *
FROM computer
WHERE computer_id = $1;

-- name: CreateComputer :one
INSERT INTO computer (
    hostname,
    cpu_id,
    ram_gb,
    ram_type,
    storage_gb,
    storage_type,
    computer_type,
    observations,
    equipment_user_id,
    room_id,
    mac_address,
    created_by_app_user_id
) VALUES (
    @hostname,
    sqlc.narg(cpu_id),
    @ram_gb,
    @ram_type,
    @storage_gb,
    @storage_type,
    @computer_type,
    sqlc.narg(observations),
    sqlc.narg(equipment_user_id),
    sqlc.narg(room_id),
    sqlc.narg(mac_address),
    @created_by_app_user_id
)
RETURNING *;

-- name: UpdateComputer :one
UPDATE computer
SET
    hostname         = COALESCE(sqlc.narg(hostname),         hostname),
    cpu_id           = COALESCE(sqlc.narg(cpu_id),           cpu_id),
    ram_gb           = COALESCE(sqlc.narg(ram_gb),           ram_gb),
    ram_type         = COALESCE(sqlc.narg(ram_type),         ram_type),
    storage_gb       = COALESCE(sqlc.narg(storage_gb),       storage_gb),
    storage_type     = COALESCE(sqlc.narg(storage_type),     storage_type),
    computer_type    = COALESCE(sqlc.narg(computer_type),    computer_type),
    observations     = COALESCE(sqlc.narg(observations),     observations),
    equipment_user_id = COALESCE(sqlc.narg(equipment_user_id), equipment_user_id),
    room_id          = COALESCE(sqlc.narg(room_id),          room_id),
    mac_address      = COALESCE(sqlc.narg(mac_address),      mac_address),
    updated_at       = now()
WHERE computer_id = sqlc.arg(computer_id)
RETURNING *;

-- name: DeleteComputer :exec
DELETE FROM computer
WHERE computer_id = $1;

-- name: ListComputerOS :many
SELECT o.*
FROM os o
INNER JOIN computer_os co ON co.os_id = o.os_id
WHERE co.computer_id = $1
ORDER BY o.os_id;

-- name: AddComputerOS :exec
INSERT INTO computer_os (computer_id, os_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveComputerOS :exec
DELETE FROM computer_os
WHERE computer_id = $1 AND os_id = $2;

-- name: InsertComputerAudit :exec
INSERT INTO computer_audit (
    event_type,
    computer_id,
    old_values,
    new_values,
    changed_by_app_user_id
) VALUES (
    @event_type,
    @computer_id,
    @old_values,
    sqlc.narg(new_values),
    @changed_by_app_user_id
);

-- name: GetComputerAudit :many
SELECT *
FROM computer_audit
WHERE computer_id = $1
ORDER BY changed_at DESC;
