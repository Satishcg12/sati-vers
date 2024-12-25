-- name: CreateSalt :one
INSERT INTO salts (user_id, salt_value)
VALUES ($1, $2)
RETURNING id;
