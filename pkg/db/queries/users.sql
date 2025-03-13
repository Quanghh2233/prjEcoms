-- name: CreateUser :one
INSERT INTO users (username, email, password_hash, role)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at
LIMIT $1 OFFSET $2;

-- name: UpdateUserRole :one
UPDATE users
SET role = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: BanUser :one
UPDATE users
SET is_banned = true, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UnbanUser :one
UPDATE users
SET is_banned = false, updated_at = now()
WHERE id = $1
RETURNING *;
