package token

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWT struct {
	secret string
}

func NewJWT(secret string) Token {
	return &JWT{
		secret: secret,
	}
}

func (j *JWT) Create(sub string, exp time.Duration) (*Response, error) {
	payload := NewPayload(sub, exp)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tk, err := t.SignedString([]byte(j.secret))
	if err != nil {
		return nil, err
	}

	rsp := &Response{
		ExpiresIn: payload.ExpireIn,
		Unit:      "ms",
		Type:      "bearer",
		Token:     tk,
	}
	return rsp, nil
}
