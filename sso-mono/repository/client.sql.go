// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: client.sql

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const getClientById = `-- name: GetClientById :one
SELECT client_id, owner_id, client_secret_hash, client_name, client_description, clinet_logo_url, client_tos_url, client_policy_url, client_homepage_url, client_redirect_uris, client_scopes, client_grants, is_trusted, status, created_by, created_at, updated_at
FROM clients
WHERE client_id = $1
LIMIT 1
`

func (q *Queries) GetClientById(ctx context.Context, clientID uuid.UUID) (Client, error) {
	row := q.db.QueryRowContext(ctx, getClientById, clientID)
	var i Client
	err := row.Scan(
		&i.ClientID,
		&i.OwnerID,
		&i.ClientSecretHash,
		&i.ClientName,
		&i.ClientDescription,
		&i.ClinetLogoUrl,
		&i.ClientTosUrl,
		&i.ClientPolicyUrl,
		&i.ClientHomepageUrl,
		pq.Array(&i.ClientRedirectUris),
		pq.Array(&i.ClientScopes),
		pq.Array(&i.ClientGrants),
		&i.IsTrusted,
		&i.Status,
		&i.CreatedBy,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
