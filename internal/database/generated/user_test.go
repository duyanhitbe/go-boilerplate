package db

import (
	"context"
	"testing"

	"github.com/duyanhitbe/go-boilerplate/internal/utils"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) *User {
	user, err := testQueries.CreateUser(context.Background(), CreateUserParams{
		Username: utils.RandomString(6),
		Password: utils.RandomString(6),
	})

	require.NoError(t, err)
	require.NotZero(t, user.ID)
	require.NotZero(t, user.Username)
	require.NotZero(t, user.Password)
	require.NotZero(t, user.CreatedAt)
	require.NotZero(t, user.UpdatedAt)

	return user
}

func TestCreate(t *testing.T) {
	createRandomUser(t)
}

func TestFindOneUserById(t *testing.T) {
	u := createRandomUser(t)

	user, err := testQueries.FindOneUserById(context.Background(), u.ID)

	require.NoError(t, err)
	require.Equal(t, u.ID, user.ID)
	require.Equal(t, u.Username, user.Username)
	require.Equal(t, u.Password, user.Password)
	require.Equal(t, u.CreatedAt, user.CreatedAt)
	require.Equal(t, u.UpdatedAt, user.UpdatedAt)
}

func TestFindOneUserByUsername(t *testing.T) {
	u := createRandomUser(t)

	user, err := testQueries.FindOneUserByUsername(context.Background(), u.Username)

	require.NoError(t, err)
	require.Equal(t, u.ID, user.ID)
	require.Equal(t, u.Username, user.Username)
	require.Equal(t, u.Password, user.Password)
	require.Equal(t, u.CreatedAt, user.CreatedAt)
	require.Equal(t, u.UpdatedAt, user.UpdatedAt)
}
