// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING user_id, username, email, password_hash, mfa_secret, mfa_enabled, is_verified, status, created_at, updated_at
`

type CreateUserParams struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.MfaSecret,
		&i.MfaEnabled,
		&i.IsVerified,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, username, email, password_hash, mfa_secret, mfa_enabled, is_verified, status, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.MfaSecret,
		&i.MfaEnabled,
		&i.IsVerified,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT user_id, username, email, password_hash, mfa_secret, mfa_enabled, is_verified, status, created_at, updated_at
FROM users
WHERE user_id = $1
LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, userID uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, userID)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.MfaSecret,
		&i.MfaEnabled,
		&i.IsVerified,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT user_id, username, email, password_hash, mfa_secret, mfa_enabled, is_verified, status, created_at, updated_at
FROM users
WHERE username = $1
LIMIT 1
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.MfaSecret,
		&i.MfaEnabled,
		&i.IsVerified,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByUsernameOrEmail = `-- name: GetUserByUsernameOrEmail :one
SELECT user_id, username, email, password_hash, mfa_secret, mfa_enabled, is_verified, status, created_at, updated_at
FROM users
WHERE username = $1 OR email = $1
LIMIT 1
`

func (q *Queries) GetUserByUsernameOrEmail(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByUsernameOrEmail, username)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.MfaSecret,
		&i.MfaEnabled,
		&i.IsVerified,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT user_id, username, email, password_hash, mfa_secret, mfa_enabled, is_verified, status, created_at, updated_at 
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type GetUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetUsers(ctx context.Context, arg GetUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.Email,
			&i.PasswordHash,
			&i.MfaSecret,
			&i.MfaEnabled,
			&i.IsVerified,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
