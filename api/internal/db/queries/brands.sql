-- name: ListBrands :many
SELECT * FROM brand ORDER BY name;

-- name: GetBrand :one
SELECT * FROM brand WHERE brand_id = $1;

-- name: CreateBrand :one
INSERT INTO brand (name) VALUES ($1) RETURNING *;

-- name: UpdateBrand :one
UPDATE brand SET name = $2 WHERE brand_id = $1 RETURNING *;

-- name: DeleteBrand :exec
DELETE FROM brand WHERE brand_id = $1;
