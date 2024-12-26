-- name: GetSaltByUserId :one
SELECT *
FROM salts
WHERE user_id = $1;

-- name: CreateSalt :one
INSERT INTO salts (user_id, salt_value)
VALUES ($1, $2)
RETURNING id;
