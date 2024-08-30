package configs

import (
	"os"
	"strconv"
)

type Env struct {
	Port  int
	DbUrl string
}

func InitEnv() *Env {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbUrl := os.Getenv("DB_URL")

	return &Env{
		Port:  port,
		DbUrl: dbUrl,
	}
}
