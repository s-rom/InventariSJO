-- ============================================================
-- PRINTER MODELS
-- ============================================================

-- name: ListPrinterModels :many
SELECT pm.*, b.name AS brand_name
FROM printer_model pm
JOIN brand b ON b.brand_id = pm.brand_id
ORDER BY b.name, pm.model_name;

-- name: GetPrinterModel :one
SELECT pm.*, b.name AS brand_name
FROM printer_model pm
JOIN brand b ON b.brand_id = pm.brand_id
WHERE pm.printer_model_id = $1;

-- name: CreatePrinterModel :one
INSERT INTO printer_model (brand_id, model_name, printer_type, print_color)
VALUES (@brand_id, @model_name, @printer_type, @print_color)
RETURNING *;

-- name: UpdatePrinterModel :one
UPDATE printer_model
SET
    brand_id     = COALESCE(sqlc.narg(brand_id),     brand_id),
    model_name   = COALESCE(sqlc.narg(model_name),   model_name),
    printer_type = COALESCE(sqlc.narg(printer_type), printer_type),
    print_color  = COALESCE(sqlc.narg(print_color),  print_color)
WHERE printer_model_id = sqlc.arg(printer_model_id)
RETURNING *;

-- name: DeletePrinterModel :exec
DELETE FROM printer_model WHERE printer_model_id = $1;


-- ============================================================
-- PRINTER SUPPLIES (consumibles)
-- ============================================================

-- name: ListPrinterSupplies :many
SELECT * FROM printer_supply ORDER BY name;

-- name: GetPrinterSupply :one
SELECT * FROM printer_supply WHERE printer_supply_id = $1;

-- name: ListSuppliesByPrinterModel :many
SELECT ps.*
FROM printer_supply ps
JOIN printer_model_supply pms ON pms.printer_supply_id = ps.printer_supply_id
WHERE pms.printer_model_id = $1
ORDER BY ps.name;

-- name: CreatePrinterSupply :one
INSERT INTO printer_supply (name, supply_type)
VALUES (@name, @supply_type)
RETURNING *;

-- name: UpdatePrinterSupply :one
UPDATE printer_supply
SET
    name        = COALESCE(sqlc.narg(name),        name),
    supply_type = COALESCE(sqlc.narg(supply_type), supply_type)
WHERE printer_supply_id = sqlc.arg(printer_supply_id)
RETURNING *;

-- name: DeletePrinterSupply :exec
DELETE FROM printer_supply WHERE printer_supply_id = $1;

-- name: AddPrinterModelSupply :exec
INSERT INTO printer_model_supply (printer_model_id, printer_supply_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemovePrinterModelSupply :exec
DELETE FROM printer_model_supply
WHERE printer_model_id = $1 AND printer_supply_id = $2;


-- ============================================================
-- PRINTERS (unitats)
-- ============================================================

-- name: ListPrinters :many
SELECT
    p.printer_id, p.status,
    p.has_network_capability, p.uses_network, p.ip_address,
    p.room_id, p.equipment_user_id, p.observations,
    p.created_by_app_user_id, p.created_at, p.updated_at,
    p.printer_model_id,
    pm.model_name,
    pm.printer_type,
    pm.print_color,
    b.brand_id,
    b.name AS brand_name,
    r.name AS room_name,
    eu.name AS equipment_user_name
FROM printer p
JOIN printer_model pm ON pm.printer_model_id = p.printer_model_id
JOIN brand b ON b.brand_id = pm.brand_id
LEFT JOIN room r ON r.room_id = p.room_id
LEFT JOIN equipment_user eu ON eu.equipment_user_id = p.equipment_user_id
ORDER BY p.printer_id;

-- name: GetPrinter :one
SELECT
    p.printer_id, p.status,
    p.has_network_capability, p.uses_network, p.ip_address,
    p.room_id, p.equipment_user_id, p.observations,
    p.created_by_app_user_id, p.created_at, p.updated_at,
    p.printer_model_id,
    pm.model_name,
    pm.printer_type,
    pm.print_color,
    b.brand_id,
    b.name AS brand_name,
    r.name AS room_name,
    eu.name AS equipment_user_name
FROM printer p
JOIN printer_model pm ON pm.printer_model_id = p.printer_model_id
JOIN brand b ON b.brand_id = pm.brand_id
LEFT JOIN room r ON r.room_id = p.room_id
LEFT JOIN equipment_user eu ON eu.equipment_user_id = p.equipment_user_id
WHERE p.printer_id = $1;

-- name: CreatePrinter :one
INSERT INTO printer (
    printer_model_id,
    status,
    has_network_capability,
    uses_network,
    ip_address,
    room_id,
    equipment_user_id,
    observations,
    created_by_app_user_id
) VALUES (
    @printer_model_id,
    @status,
    @has_network_capability,
    @uses_network,
    sqlc.narg(ip_address),
    sqlc.narg(room_id),
    sqlc.narg(equipment_user_id),
    sqlc.narg(observations),
    @created_by_app_user_id
) RETURNING *;

-- name: UpdatePrinter :one
UPDATE printer
SET
    printer_model_id       = COALESCE(sqlc.narg(printer_model_id),       printer_model_id),
    status                 = COALESCE(sqlc.narg(status),                 status),
    has_network_capability = COALESCE(sqlc.narg(has_network_capability), has_network_capability),
    uses_network           = COALESCE(sqlc.narg(uses_network),           uses_network),
    ip_address             = COALESCE(sqlc.narg(ip_address),             ip_address),
    room_id                = COALESCE(sqlc.narg(room_id),                room_id),
    equipment_user_id      = COALESCE(sqlc.narg(equipment_user_id),      equipment_user_id),
    observations           = COALESCE(sqlc.narg(observations),           observations),
    updated_at             = now()
WHERE printer_id = sqlc.arg(printer_id)
RETURNING *;

-- name: DeletePrinter :exec
DELETE FROM printer WHERE printer_id = $1;
