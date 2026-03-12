-- name: ListDesktopModels :many
SELECT dm.*, b.name AS brand_name
FROM desktop_model dm
JOIN brand b ON b.brand_id = dm.brand_id
ORDER BY b.name, dm.model_name;

-- name: GetDesktopModel :one
SELECT dm.*, b.name AS brand_name
FROM desktop_model dm
JOIN brand b ON b.brand_id = dm.brand_id
WHERE dm.desktop_model_id = $1;

-- name: CreateDesktopModel :one
INSERT INTO desktop_model (
    brand_id, model_name, cpu_id,
    base_ram_gb, base_ram_type,
    base_storage_gb, base_storage_type, base_os_id
) VALUES (
    @brand_id, @model_name, sqlc.narg(cpu_id),
    @base_ram_gb, @base_ram_type,
    @base_storage_gb, @base_storage_type, sqlc.narg(base_os_id)
) RETURNING *;

-- name: UpdateDesktopModel :one
UPDATE desktop_model
SET
    brand_id          = COALESCE(sqlc.narg(brand_id),          brand_id),
    model_name        = COALESCE(sqlc.narg(model_name),        model_name),
    cpu_id            = COALESCE(sqlc.narg(cpu_id),            cpu_id),
    base_ram_gb       = COALESCE(sqlc.narg(base_ram_gb),       base_ram_gb),
    base_ram_type     = COALESCE(sqlc.narg(base_ram_type),     base_ram_type),
    base_storage_gb   = COALESCE(sqlc.narg(base_storage_gb),   base_storage_gb),
    base_storage_type = COALESCE(sqlc.narg(base_storage_type), base_storage_type),
    base_os_id        = COALESCE(sqlc.narg(base_os_id),        base_os_id)
WHERE desktop_model_id = sqlc.arg(desktop_model_id)
RETURNING *;

-- name: DeleteDesktopModel :exec
DELETE FROM desktop_model WHERE desktop_model_id = $1;
