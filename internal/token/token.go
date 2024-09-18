package token

import "time"

type Response struct {
	ExpiresIn int64
	Unit      string
	Type      string
	Token     string
}

type Token interface {
	Create(sub string, exp time.Duration) (*Response, error)
}
