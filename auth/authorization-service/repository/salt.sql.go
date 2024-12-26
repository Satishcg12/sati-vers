// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: salt.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const createSalt = `-- name: CreateSalt :one
INSERT INTO salts (user_id, salt_value)
VALUES ($1, $2)
RETURNING id
`

type CreateSaltParams struct {
	UserID    uuid.UUID `json:"user_id"`
	SaltValue string    `json:"salt_value"`
}

func (q *Queries) CreateSalt(ctx context.Context, arg CreateSaltParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createSalt, arg.UserID, arg.SaltValue)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getSaltByUserId = `-- name: GetSaltByUserId :one
SELECT id, user_id, salt_value, created_at, updated_at
FROM salts
WHERE user_id = $1
`

func (q *Queries) GetSaltByUserId(ctx context.Context, userID uuid.UUID) (Salt, error) {
	row := q.db.QueryRowContext(ctx, getSaltByUserId, userID)
	var i Salt
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SaltValue,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
