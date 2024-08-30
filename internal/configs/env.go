package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct {
	Port  int
	DbUrl string
}

func InitEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	dbUrl := os.Getenv("DB_URL")

	return &Env{
		Port:  port,
		DbUrl: dbUrl,
	}
}
