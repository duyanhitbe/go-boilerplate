package token

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

const iss = "https://github.com/duyanhitbe/go-boilerplate"

type Payload struct {
	jwt.RegisteredClaims
	ExpireIn int64 `json:"expire_in"`
}

func NewPayload(sub string, exp time.Duration) *Payload {
	now := time.Now()
	id, _ := uuid.NewUUID()
	return &Payload{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			Subject:   sub,
			Issuer:    iss,
			Audience:  []string{iss},
			ExpiresAt: jwt.NewNumericDate(now.Add(exp)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		ExpireIn: exp.Milliseconds(),
	}
}
