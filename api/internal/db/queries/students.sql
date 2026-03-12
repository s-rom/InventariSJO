-- name: ListStudentsByClass :many
SELECT * FROM student WHERE class_id = $1 ORDER BY full_name;

-- name: GetStudent :one
SELECT * FROM student WHERE student_id = $1;

-- name: CreateStudent :one
INSERT INTO student (full_name, class_id)
VALUES (@full_name, @class_id)
RETURNING *;

-- name: UpdateStudent :one
UPDATE student
SET
    full_name = COALESCE(sqlc.narg(full_name), full_name),
    class_id  = COALESCE(sqlc.narg(class_id),  class_id)
WHERE student_id = sqlc.arg(student_id)
RETURNING *;

-- name: DeleteStudent :exec
DELETE FROM student WHERE student_id = $1;
