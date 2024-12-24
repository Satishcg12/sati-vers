// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID      `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	CreatedAt sql.NullTime   `json:"created_at"`
	UpdatedAt sql.NullTime   `json:"updated_at"`
	Status    sql.NullString `json:"status"`
}

type UserCredential struct {
	ID           uuid.UUID     `json:"id"`
	UserID       uuid.NullUUID `json:"user_id"`
	PasswordHash string        `json:"password_hash"`
	CreatedAt    sql.NullTime  `json:"created_at"`
	UpdatedAt    sql.NullTime  `json:"updated_at"`
}
