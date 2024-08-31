package main

import (
	"database/sql"
	"log"

	"github.com/duyanhitbe/go-boilerplate/internal/configs"
	db "github.com/duyanhitbe/go-boilerplate/internal/database/generated"
	"github.com/duyanhitbe/go-boilerplate/internal/server"
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
	conn, err := sql.Open("postgres", env.DbUrl)
	if err != nil {
		log.Fatal(err)
	}

	//create sql store
	store := db.NewSQLStore(conn)

	//init API server
	server := server.NewServer(env.Port, store)
	if err = server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
