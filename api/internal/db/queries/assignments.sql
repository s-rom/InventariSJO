-- name: ListAssignmentsByLaptop :many
SELECT
    lsa.assignment_id, lsa.computer_id, lsa.student_id,
    lsa.class_id, lsa.academic_year,
    s.full_name AS student_name,
    sc.course, sc.class_label, sc.shift, sc.cycle_id
FROM laptop_student_assignment lsa
JOIN student s   ON s.student_id     = lsa.student_id
JOIN school_class sc ON sc.class_id  = lsa.class_id
WHERE lsa.computer_id = $1
ORDER BY lsa.academic_year DESC, sc.shift;

-- name: GetAssignment :one
SELECT * FROM laptop_student_assignment WHERE assignment_id = $1;

-- name: CreateAssignment :one
INSERT INTO laptop_student_assignment (computer_id, student_id, class_id, academic_year)
VALUES (@computer_id, @student_id, @class_id, @academic_year)
RETURNING *;

-- name: UpdateAssignment :one
UPDATE laptop_student_assignment
SET
    student_id    = COALESCE(sqlc.narg(student_id),    student_id),
    class_id      = COALESCE(sqlc.narg(class_id),      class_id),
    academic_year = COALESCE(sqlc.narg(academic_year), academic_year)
WHERE assignment_id = sqlc.arg(assignment_id)
RETURNING *;

-- name: DeleteAssignment :exec
DELETE FROM laptop_student_assignment WHERE assignment_id = $1;

-- name: ListAssignmentsByClass :many
SELECT
    lsa.assignment_id, lsa.computer_id, lsa.student_id,
    lsa.class_id, lsa.academic_year,
    s.full_name AS student_name
FROM laptop_student_assignment lsa
JOIN student s ON s.student_id = lsa.student_id
WHERE lsa.class_id = $1
ORDER BY lsa.academic_year DESC, s.full_name;
