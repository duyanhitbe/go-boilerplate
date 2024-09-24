package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	s := RandomString(6)

	require.NotEmpty(t, s)
	require.Equal(t, len(s), 6)
}

func TestRandomInt(t *testing.T) {
	i := RandomInt(10)

	require.GreaterOrEqual(t, i, 0)
	require.Less(t, i, 10)
}
