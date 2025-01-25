-- name: CreateAuthCode :one
INSERT INTO auth_codes (client_id, user_id, auth_code_hash, redirect_uri, scopes, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
