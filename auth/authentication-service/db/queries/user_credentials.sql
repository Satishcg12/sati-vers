-- name: GetCredentialsByUserId :one
SELECT *
FROM user_credentials
WHERE user_id = $1;

-- name: CreateCredentials :one
INSERT INTO user_credentials (user_id, credential_type, credential_value)
VALUES ($1, $2, $3)
RETURNING *;