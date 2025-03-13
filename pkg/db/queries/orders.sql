-- name: CreateOrder :one
INSERT INTO orders (user_id, total_amount, shipping_address)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AddOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_at_purchase)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: GetOrderStatistics :one
SELECT 
    COUNT(*) as total_orders,
    SUM(total_amount) as total_revenue
FROM orders
WHERE shop_id = $1 AND created_at >= $2 AND created_at <= $3;
