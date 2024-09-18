package token

import (
	"errors"
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

func (j *JWT) Verify(token string) (*Payload, error) {
	var payload Payload

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	}

	t, err := jwt.ParseWithClaims(token, &payload, keyFunc)
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	return &payload, nil
}
