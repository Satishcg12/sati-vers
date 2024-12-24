-- name: GetCredentialsByUserId :one
SELECT *
FROM user_credentials
WHERE user_id = $1;

-- name: CreateCredentials :one
INSERT INTO user_credentials (user_id, password_hash)
VALUES ($1, $2)
RETURNING *;
