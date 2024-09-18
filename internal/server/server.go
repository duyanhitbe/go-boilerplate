package server

import (
	"fmt"
	"github.com/duyanhitbe/go-boilerplate/internal/hash"
	"github.com/duyanhitbe/go-boilerplate/internal/token"
	"net/http"

	db "github.com/duyanhitbe/go-boilerplate/internal/database/generated"
)

type Server struct {
	store db.Store
	h     hash.Hash
	t     token.Token
}

func NewServer(addr string, store db.Store, h hash.Hash, t token.Token) *http.Server {
	s := &Server{
		store: store,
		h:     h,
		t:     t,
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", addr),
		Handler: s.router(),
	}
}
