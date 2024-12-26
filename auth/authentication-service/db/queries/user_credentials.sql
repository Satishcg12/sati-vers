-- name: GetCredentialsByUserId :one
SELECT *
FROM user_credentials
WHERE user_id = $1;


-- name: GetSaltAndCredentialsByUserId :one
SELECT sqlc.embed(user_credentials), sqlc.embed(salts)
FROM user_credentials
JOIN salts ON user_credentials.user_id = salts.user_id
WHERE user_credentials.user_id = $1;

-- name: CreateCredentials :one
INSERT INTO user_credentials (user_id, credential_type, credential_value)
VALUES ($1, $2, $3)
RETURNING *;