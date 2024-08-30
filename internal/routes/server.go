package routes

import (
	"database/sql"
	"fmt"

	db "github.com/duyanhitbe/go-boilerplate/internal/database/generated"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store db.Store
}

func NewServer(database *sql.DB) *Server {
	return &Server{
		store: db.NewSQLStore(database),
	}
}

func (s *Server) Start(port int) error {
	r := gin.Default()

	s.registerIndexRoutes(r)
	s.registerTodoRoutes(r)

	return r.Run(fmt.Sprintf(":%d", port))
}
