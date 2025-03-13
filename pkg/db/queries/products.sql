-- name: CreateProduct :one
INSERT INTO products (shop_id, name, description, price, stock, category, image_urls)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1;

-- name: GetProductsByShopID :many
SELECT * FROM products
WHERE shop_id = $1
ORDER BY created_at
LIMIT $2 OFFSET $3;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY created_at
LIMIT $1 OFFSET $2;

-- name: SearchProducts :many
SELECT * FROM products
WHERE name ILIKE $1 OR description ILIKE $1
ORDER BY created_at
LIMIT $2 OFFSET $3;

-- name: FilterProductsByCategory :many
SELECT * FROM products
WHERE category = $1
ORDER BY created_at
LIMIT $2 OFFSET $3;

-- name: FilterProductsByPrice :many
SELECT * FROM products
WHERE price BETWEEN $1 AND $2
ORDER BY created_at
LIMIT $3 OFFSET $4;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, price = $4, 
    stock = $5, category = $6, image_urls = $7, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;
