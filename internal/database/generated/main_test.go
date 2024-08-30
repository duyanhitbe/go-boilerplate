package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/duyanhitbe/go-boilerplate/internal/configs"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	env := configs.InitEnv()

	testDB, err := sql.Open("postgres", env.DbUrl)
	if err != nil {
		log.Fatal("Can not connect database", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
