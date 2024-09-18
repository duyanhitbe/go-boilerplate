// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (*User, error)
	FetchUser(ctx context.Context, arg FetchUserParams) ([]*User, error)
	FindOneUserById(ctx context.Context, id uuid.UUID) (*User, error)
	FindOneUserByUsername(ctx context.Context, username string) (*User, error)
}

var _ Querier = (*Queries)(nil)
