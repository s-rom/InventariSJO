-- name: ListCPUs :many
SELECT *
FROM cpu
ORDER BY cpu_id;

-- name: CreateCPU :one
INSERT INTO cpu (model_name, benchmark_score)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateCPU :one
UPDATE cpu
SET
    model_name      = COALESCE(sqlc.narg(model_name),      model_name),
    benchmark_score = COALESCE(sqlc.narg(benchmark_score), benchmark_score)
WHERE cpu_id = sqlc.arg(cpu_id)
RETURNING *;

-- name: DeleteCPU :exec
DELETE FROM cpu
WHERE cpu_id = $1;
