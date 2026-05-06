-- ============================================================
-- PROJECTOR MODELS
-- ============================================================

-- name: ListProjectorModels :many
SELECT pm.*, b.name AS brand_name
FROM projector_model pm
JOIN brand b ON b.brand_id = pm.brand_id
ORDER BY b.name, pm.model_name;

-- name: GetProjectorModel :one
SELECT pm.*, b.name AS brand_name
FROM projector_model pm
JOIN brand b ON b.brand_id = pm.brand_id
WHERE pm.projector_model_id = $1;

-- name: CreateProjectorModel :one
INSERT INTO projector_model (brand_id, model_name)
VALUES (@brand_id, @model_name)
RETURNING *;

-- name: UpdateProjectorModel :one
UPDATE projector_model
SET
    brand_id   = COALESCE(sqlc.narg(brand_id),   brand_id),
    model_name = COALESCE(sqlc.narg(model_name), model_name)
WHERE projector_model_id = sqlc.arg(projector_model_id)
RETURNING *;

-- name: DeleteProjectorModel :exec
DELETE FROM projector_model WHERE projector_model_id = $1;


-- ============================================================
-- PROJECTORS (unitats)
-- ============================================================

-- name: ListProjectors :many
SELECT
    p.projector_id, p.serial_number, p.status,
    p.room_id, p.equipment_user_id, p.observations,
    p.created_by_app_user_id, p.created_at, p.updated_at,
    p.projector_model_id,
    pm.model_name,
    pm.brand_id,
    b.name AS brand_name,
    r.name AS room_name,
    eu.name AS equipment_user_name
FROM projector p
JOIN projector_model pm ON pm.projector_model_id = p.projector_model_id
JOIN brand b ON b.brand_id = pm.brand_id
LEFT JOIN room r ON r.room_id = p.room_id
LEFT JOIN equipment_user eu ON eu.equipment_user_id = p.equipment_user_id
ORDER BY p.projector_id;

-- name: GetProjector :one
SELECT
    p.projector_id, p.serial_number, p.status,
    p.room_id, p.equipment_user_id, p.observations,
    p.created_by_app_user_id, p.created_at, p.updated_at,
    p.projector_model_id,
    pm.model_name,
    pm.brand_id,
    b.name AS brand_name,
    r.name AS room_name,
    eu.name AS equipment_user_name
FROM projector p
JOIN projector_model pm ON pm.projector_model_id = p.projector_model_id
JOIN brand b ON b.brand_id = pm.brand_id
LEFT JOIN room r ON r.room_id = p.room_id
LEFT JOIN equipment_user eu ON eu.equipment_user_id = p.equipment_user_id
WHERE p.projector_id = $1;

-- name: CreateProjector :one
INSERT INTO projector (
    projector_model_id,
    serial_number,
    status,
    room_id,
    equipment_user_id,
    observations,
    created_by_app_user_id
) VALUES (
    @projector_model_id,
    sqlc.narg(serial_number),
    @status,
    sqlc.narg(room_id),
    sqlc.narg(equipment_user_id),
    sqlc.narg(observations),
    @created_by_app_user_id
) RETURNING *;

-- name: UpdateProjector :one
UPDATE projector
SET
    projector_model_id = COALESCE(sqlc.narg(projector_model_id), projector_model_id),
    serial_number      = COALESCE(sqlc.narg(serial_number),      serial_number),
    status             = COALESCE(sqlc.narg(status),             status),
    room_id            = COALESCE(sqlc.narg(room_id),            room_id),
    equipment_user_id  = COALESCE(sqlc.narg(equipment_user_id),  equipment_user_id),
    observations       = sqlc.narg(observations),
    updated_at         = now()
WHERE projector_id = sqlc.arg(projector_id)
RETURNING *;

-- name: DeleteProjector :exec
DELETE FROM projector WHERE projector_id = $1;
