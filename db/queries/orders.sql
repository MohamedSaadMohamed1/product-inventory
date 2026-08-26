-- name: CreateOrder :one
INSERT INTO orders (status, total_amount, user_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: GetOrderForUpdate :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1
FOR UPDATE;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, cancelled_at = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, unit_price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListOrderItemsByOrderId :many
SELECT * FROM order_items
WHERE order_id = $1
ORDER BY id ASC;

-- name: ListOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY id DESC
LIMIT $2 OFFSET $3;
