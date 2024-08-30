package main

import (
	"database/sql"
	"log"

	"github.com/duyanhitbe/go-boilerplate/internal/configs"
	"github.com/duyanhitbe/go-boilerplate/internal/routes"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	//init environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	env := configs.InitEnv()

	//open db connection
	db, err := sql.Open("postgres", env.DbUrl)
	if err != nil {
		log.Fatal(err)
	}

	//init API server
	server := routes.NewServer(db)
	if err = server.Start(env.Port); err != nil {
		log.Fatal(err)
	}
}
