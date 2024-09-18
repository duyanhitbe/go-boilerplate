package main

import (
	"database/sql"
	"github.com/duyanhitbe/go-boilerplate/internal/hash"
	"github.com/duyanhitbe/go-boilerplate/internal/token"
	"github.com/gin-gonic/gin"
	"os"

	"github.com/duyanhitbe/go-boilerplate/internal/configs"
	db "github.com/duyanhitbe/go-boilerplate/internal/database/generated"
	"github.com/duyanhitbe/go-boilerplate/internal/server"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	gin.SetMode(gin.ReleaseMode)

	//init environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
	log.Info().Msg("Loaded env")
	env := configs.InitEnv()

	//open db connection
	conn, err := sql.Open("postgres", env.DbUrl)
	if err != nil {
		log.Fatal().Err(err)
	}
	log.Info().Msg("Connected to database")

	//create sql store
	store := db.NewSQLStore(conn)

	//create hash
	h := hash.NewBcrypt()

	//create token
	t := token.NewJWT("secret")

	//init API server
	s := server.NewServer(env.Port, store, h, t)
	log.Info().Msg("Started server")
	log.Info().Msg("Server is running on port " + env.Port)
	if err = s.ListenAndServe(); err != nil {
		log.Fatal().Err(err)
	}
}
