package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	title = "title test"
)

func createRandomTodo(t *testing.T) *Todo {
	data, err := testQueries.CreateTodo(context.Background(), title)

	require.NoError(t, err)
	require.NotNil(t, data)
	require.NotZero(t, data)

	return data
}

func removeAllTodo(t *testing.T) {
	err := testQueries.DeleteAllTodo(context.Background())
	require.NoError(t, err)
}

func TestGetAllTodo(t *testing.T) {
	removeAllTodo(t)
	createRandomTodo(t)

	data, err := testQueries.GetAllTodo(context.Background())

	require.NoError(t, err)
	require.NotNil(t, data)
	require.NotZero(t, data)
	require.Equal(t, data[0].Title, title)

	removeAllTodo(t)
}
