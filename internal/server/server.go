package server

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/duyanhitbe/go-boilerplate/internal/hash"
	"github.com/duyanhitbe/go-boilerplate/internal/token"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

func (s *Server) getUserFromContext(c *gin.Context) (*db.User, bool) {
	authPayload, valid := c.Get(authenticationPayload)
	if !valid {
		log.Error().Msg("Can not get payload from context")
		throwInternalServerError(c, errors.New("invalid payload"))
		return nil, false
	}

	payload, valid := authPayload.(*token.Payload)
	if !valid {
		log.Error().Msg("Can not parse payload to token.Payload")
		throwInternalServerError(c, errors.New("invalid payload"))
		return nil, false
	}

	id, err := uuid.Parse(payload.Subject)
	if err != nil {
		log.Error().Msg("Can not parse subject to uuid")
		throwInternalServerError(c, errors.New("invalid payload"))
		return nil, false
	}

	user, err := s.store.FindOneUserById(c, id)
	if err != nil {
		log.Error().Msgf("FindOneUserById fail %v", err)

		if errors.Is(err, sql.ErrNoRows) {
			log.Error().Msg("user not found")
			throwForbidden(c, errors.New("user not found"))
			return nil, false
		}

		throwInternalServerError(c, err)
		return nil, false
	}

	return user, true
}
