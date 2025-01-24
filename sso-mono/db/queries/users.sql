-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUsers :many
SELECT * 
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByUsernameOrEmail :one
SELECT *
FROM users
WHERE username = $1 OR email = $1
LIMIT 1;
    

-- name: GetUserById :one
SELECT *
FROM users
WHERE user_id = $1
LIMIT 1;

