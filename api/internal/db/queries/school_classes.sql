-- name: ListClasses :many
SELECT sc.*, cy.name AS cycle_name
FROM school_class sc
JOIN cycle cy ON cy.cycle_id = sc.cycle_id
ORDER BY cy.name, sc.course, sc.class_label;

-- name: ListClassesByCycle :many
SELECT sc.*, cy.name AS cycle_name
FROM school_class sc
JOIN cycle cy ON cy.cycle_id = sc.cycle_id
WHERE sc.cycle_id = $1
ORDER BY sc.course, sc.class_label;

-- name: GetClass :one
SELECT * FROM school_class WHERE class_id = $1;

-- name: CreateClass :one
INSERT INTO school_class (cycle_id, course, class_label, shift, tutor_app_user_id)
VALUES (@cycle_id, @course, @class_label, @shift, sqlc.narg(tutor_app_user_id))
RETURNING *;

-- name: UpdateClass :one
UPDATE school_class
SET
    class_label       = COALESCE(sqlc.narg(class_label),       class_label),
    shift             = COALESCE(sqlc.narg(shift),             shift),
    tutor_app_user_id = COALESCE(sqlc.narg(tutor_app_user_id), tutor_app_user_id)
WHERE class_id = sqlc.arg(class_id)
RETURNING *;

-- name: DeleteClass :exec
DELETE FROM school_class WHERE class_id = $1;

-- name: ListClassesByTutor :many
SELECT sc.*, cy.name AS cycle_name
FROM school_class sc
JOIN cycle cy ON cy.cycle_id = sc.cycle_id
WHERE sc.tutor_app_user_id = $1
ORDER BY cy.name, sc.course, sc.class_label;
