-- name: GetClientById :one
SELECT *
FROM clients
WHERE client_id = $1
LIMIT 1;
