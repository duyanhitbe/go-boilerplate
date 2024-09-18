package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	tokenType             = "bearer"
	authenticationPayload = "authentication_payload"
)

func (s *Server) authenticationMiddleware(c *gin.Context) {
	authorization := c.Request.Header.Get("Authorization")
	if authorization == "" {
		throwUnauthorized(c, errors.New("missing authorization"))
		return
	}

	split := strings.Split(authorization, " ")
	if len(split) != 2 || strings.ToLower(split[0]) != tokenType {
		throwUnauthorized(c, errors.New("invalid authorization"))
		return
	}

	payload, err := s.t.Verify(split[1])
	if err != nil {
		throwUnauthorized(c, err)
		return
	}

	c.Set(authenticationPayload, payload)
	c.Next()
}
