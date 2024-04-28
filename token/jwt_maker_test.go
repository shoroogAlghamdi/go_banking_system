package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shoroogAlghamdi/banking_system/util"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWT(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	issuaedAt := time.Now()
	duration := time.Minute
	exiredAt := issuaedAt.Add(duration)

	jwtToken, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, jwtToken)

	payload, err := maker.VerifyToken(jwtToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)


	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuaedAt, payload.IssuesAt, time.Second)
	require.WithinDuration(t, exiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWT(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	jwtToken, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, jwtToken)

	payload, err := maker.VerifyToken(jwtToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredtoken.Error())
	require.Nil(t, payload)
}

func TestInvalidSignAlgNone(t * testing.T) {
	payload, err := NewPayload(util.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWT(util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrorInvalidToken.Error())
	require.Nil(t, payload)
}
