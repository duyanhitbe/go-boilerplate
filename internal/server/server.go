package server

import (
	"fmt"
	"net/http"

	db "github.com/duyanhitbe/go-boilerplate/internal/database/generated"
)

type Server struct {
	store db.Store
}

func NewServer(addr string, store db.Store) *http.Server {
	s := &Server{
		store: store,
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: s.RegisterRoutes(),
	}
}
