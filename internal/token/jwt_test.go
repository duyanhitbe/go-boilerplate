package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

const (
	secret = "secret_test"
	sub    = "sub_test"
	exp    = time.Minute
)

func createToken(t *testing.T) (Token, *Response) {
	j := NewJWT(secret)

	rsp, err := j.Create(sub, exp)
	var tk string

	require.NoError(t, err)
	require.NotEmpty(t, rsp)
	require.NotZero(t, rsp)
	require.Equal(t, "bearer", rsp.Type)
	require.Equal(t, "ms", rsp.Unit)
	require.Equal(t, exp.Milliseconds(), rsp.ExpiresIn)
	require.IsType(t, tk, rsp.Token)

	return j, rsp
}

func TestCreate(t *testing.T) {
	createToken(t)
}

func TestVerify(t *testing.T) {
	j, rsp := createToken(t)

	payload, err := j.Verify(rsp.Token)
	iat := payload.IssuedAt
	expireAt := iat.Add(exp)

	require.NoError(t, err)
	require.Equal(t, jwt.ClaimStrings{iss}, payload.Audience)
	require.Equal(t, exp.Milliseconds(), payload.ExpireIn)
	require.Equal(t, iss, payload.Issuer)
	require.Equal(t, sub, payload.Subject)
	require.Equal(t, expireAt, payload.ExpiresAt.Time)
}
