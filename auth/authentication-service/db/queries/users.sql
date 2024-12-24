-- name: GetUsers :many
SELECT * 
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (username, email)
VALUES ($1, $2)
RETURNING *;
