package token

import (
	"testing"
	"time"

	"github.com/shoroogAlghamdi/banking_system/util"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPaseto(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	issuaedAt := time.Now()
	duration := time.Minute
	exiredAt := issuaedAt.Add(duration)

	pToken, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, pToken)

	payload, err := maker.VerifyToken(pToken)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuaedAt, payload.IssuesAt, time.Second)
	require.WithinDuration(t, exiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredTokenPaseto(t *testing.T) {
	maker, err := NewPaseto(util.RandomString(32))
	require.NoError(t, err)
	require.NotEmpty(t, maker)

	username := util.RandomOwner()
	pToken, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, pToken)

	payload, err := maker.VerifyToken(pToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrorExpiredtoken.Error())
	require.Nil(t, payload)
}
