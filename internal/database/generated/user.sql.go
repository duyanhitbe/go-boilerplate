// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: user.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING id, username, password, created_at, updated_at
`

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const findOneUserById = `-- name: FindOneUserById :one
SELECT id, username, password, created_at, updated_at
FROM users
WHERE id = $1
`

func (q *Queries) FindOneUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	row := q.db.QueryRowContext(ctx, findOneUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}

const findOneUserByUsername = `-- name: FindOneUserByUsername :one
SELECT id, username, password, created_at, updated_at
FROM users
WHERE username = $1
`

func (q *Queries) FindOneUserByUsername(ctx context.Context, username string) (*User, error) {
	row := q.db.QueryRowContext(ctx, findOneUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return &i, err
}
