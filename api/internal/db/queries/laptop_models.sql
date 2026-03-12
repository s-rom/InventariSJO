-- name: ListLaptopModels :many
SELECT lm.*, b.name AS brand_name
FROM laptop_model lm
JOIN brand b ON b.brand_id = lm.brand_id
ORDER BY b.name, lm.model_name;

-- name: GetLaptopModel :one
SELECT lm.*, b.name AS brand_name
FROM laptop_model lm
JOIN brand b ON b.brand_id = lm.brand_id
WHERE lm.laptop_model_id = $1;

-- name: CreateLaptopModel :one
INSERT INTO laptop_model (
    brand_id, model_name, cpu_id,
    base_ram_gb, base_ram_type,
    base_storage_gb, base_storage_type, base_os_id
) VALUES (
    @brand_id, @model_name, sqlc.narg(cpu_id),
    @base_ram_gb, @base_ram_type,
    @base_storage_gb, @base_storage_type, sqlc.narg(base_os_id)
) RETURNING *;

-- name: UpdateLaptopModel :one
UPDATE laptop_model
SET
    brand_id          = COALESCE(sqlc.narg(brand_id),          brand_id),
    model_name        = COALESCE(sqlc.narg(model_name),        model_name),
    cpu_id            = COALESCE(sqlc.narg(cpu_id),            cpu_id),
    base_ram_gb       = COALESCE(sqlc.narg(base_ram_gb),       base_ram_gb),
    base_ram_type     = COALESCE(sqlc.narg(base_ram_type),     base_ram_type),
    base_storage_gb   = COALESCE(sqlc.narg(base_storage_gb),   base_storage_gb),
    base_storage_type = COALESCE(sqlc.narg(base_storage_type), base_storage_type),
    base_os_id        = COALESCE(sqlc.narg(base_os_id),        base_os_id)
WHERE laptop_model_id = sqlc.arg(laptop_model_id)
RETURNING *;

-- name: DeleteLaptopModel :exec
DELETE FROM laptop_model WHERE laptop_model_id = $1;
