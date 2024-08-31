package configs

import (
	"os"
)

type Env struct {
	Port  string
	DbUrl string
}

func InitEnv() *Env {
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	return &Env{
		Port:  port,
		DbUrl: dbUrl,
	}
}
