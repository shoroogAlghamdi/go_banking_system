package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)
var (
	ErrorExpiredtoken error = errors.New("token has expired")
	ErrorInvalidToken error = errors.New("token is invalid")
)
type Payload struct{
	ID        uuid.UUID `json:"id"`
	Username  string `json:"username"`
	IssuesAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	paylod := *&Payload{
		ID:        tokenID,
		Username:  username,
		IssuesAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return &paylod, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrorExpiredtoken
	}
	return nil 
}