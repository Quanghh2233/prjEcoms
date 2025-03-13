-- name: CreateShop :one
INSERT INTO shops (user_id, name, description, logo_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetShopByID :one
SELECT * FROM shops
WHERE id = $1;

-- name: GetShopsByUserID :many
SELECT * FROM shops
WHERE user_id = $1;

-- name: ListShops :many
SELECT * FROM shops
ORDER BY created_at
LIMIT $1 OFFSET $2;

-- name: SearchShops :many
SELECT * FROM shops
WHERE name ILIKE $1 OR description ILIKE $1
ORDER BY created_at
LIMIT $2 OFFSET $3;

-- name: UpdateShop :one
UPDATE shops
SET name = $2, description = $3, logo_url = $4, updated_at = now()
WHERE id = $1
RETURNING *;
