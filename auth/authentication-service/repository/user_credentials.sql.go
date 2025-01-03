// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user_credentials.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createCredentials = `-- name: CreateCredentials :one
INSERT INTO user_credentials (user_id, credential_type, credential_value)
VALUES ($1, $2, $3)
RETURNING id, user_id, credential_type, credential_value, last_used, created_at, updated_at
`

type CreateCredentialsParams struct {
	UserID          uuid.NullUUID  `json:"user_id"`
	CredentialType  string         `json:"credential_type"`
	CredentialValue sql.NullString `json:"credential_value"`
}

func (q *Queries) CreateCredentials(ctx context.Context, arg CreateCredentialsParams) (UserCredential, error) {
	row := q.db.QueryRowContext(ctx, createCredentials, arg.UserID, arg.CredentialType, arg.CredentialValue)
	var i UserCredential
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CredentialType,
		&i.CredentialValue,
		&i.LastUsed,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCredentialsByUserId = `-- name: GetCredentialsByUserId :one
SELECT id, user_id, credential_type, credential_value, last_used, created_at, updated_at
FROM user_credentials
WHERE user_id = $1
`

func (q *Queries) GetCredentialsByUserId(ctx context.Context, userID uuid.NullUUID) (UserCredential, error) {
	row := q.db.QueryRowContext(ctx, getCredentialsByUserId, userID)
	var i UserCredential
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.CredentialType,
		&i.CredentialValue,
		&i.LastUsed,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSaltAndCredentialsByUserId = `-- name: GetSaltAndCredentialsByUserId :one
SELECT user_credentials.id, user_credentials.user_id, user_credentials.credential_type, user_credentials.credential_value, user_credentials.last_used, user_credentials.created_at, user_credentials.updated_at, salts.id, salts.user_id, salts.salt_value, salts.created_at, salts.updated_at
FROM user_credentials
JOIN salts ON user_credentials.user_id = salts.user_id
WHERE user_credentials.user_id = $1
`

type GetSaltAndCredentialsByUserIdRow struct {
	UserCredential UserCredential `json:"user_credential"`
	Salt           Salt           `json:"salt"`
}

func (q *Queries) GetSaltAndCredentialsByUserId(ctx context.Context, userID uuid.NullUUID) (GetSaltAndCredentialsByUserIdRow, error) {
	row := q.db.QueryRowContext(ctx, getSaltAndCredentialsByUserId, userID)
	var i GetSaltAndCredentialsByUserIdRow
	err := row.Scan(
		&i.UserCredential.ID,
		&i.UserCredential.UserID,
		&i.UserCredential.CredentialType,
		&i.UserCredential.CredentialValue,
		&i.UserCredential.LastUsed,
		&i.UserCredential.CreatedAt,
		&i.UserCredential.UpdatedAt,
		&i.Salt.ID,
		&i.Salt.UserID,
		&i.Salt.SaltValue,
		&i.Salt.CreatedAt,
		&i.Salt.UpdatedAt,
	)
	return i, err
}
