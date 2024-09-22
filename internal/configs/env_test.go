package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitEnv(t *testing.T) {
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	env := InitEnv()

	require.Equal(t, env.Port, port)
	require.Equal(t, env.DbUrl, dbUrl)
}
