package token

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/komron-dev/bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJwtMaker(t *testing.T) {
	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, payload, err := maker.CreateToken(username, role, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredJwtToken(t *testing.T) {
	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	username := util.RandomOwner()
	role := util.DepositorRole
	duration := time.Minute

	token, payload, err := maker.CreateToken(username, role, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err = maker.VerifyToken(token)
	require.Error(t, ErrExpiredToken)
	require.Nil(t, payload)
}

func TestInvalidJwtToken(t *testing.T) {
	role := util.DepositorRole
	payload, err := NewPayload(util.RandomOwner(), role, time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJwtMaker(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
