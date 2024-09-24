package hash

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	password = "test_password"
)

func createHash(t *testing.T) (Hash, string) {
	b := NewBcrypt()

	hashPwd, err := b.Create(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashPwd)
	require.NotZero(t, hashPwd)

	return b, hashPwd
}

func TestCreate(t *testing.T) {
	createHash(t)
}

func TestCompare(t *testing.T) {
	b, h := createHash(t)

	valid := b.Compare(password, h)
	notValid := b.Compare("wrong_password", h)
	require.True(t, valid)
	require.False(t, notValid)
}
