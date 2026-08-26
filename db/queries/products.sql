-- name: CreateProduct :one
INSERT INTO products (name, description, price, stock_quantity)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: GetProductForUpdate :one
SELECT * FROM products
WHERE id = $1 LIMIT 1
FOR UPDATE;

-- name: GetProductsForUpdate :many
SELECT * FROM products
WHERE id = ANY($1::bigint[])
ORDER BY id ASC
FOR UPDATE;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id ASC
LIMIT $1 OFFSET $2;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, stock_quantity = $5, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateProductStock :one
UPDATE products
SET stock_quantity = stock_quantity + $2, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
