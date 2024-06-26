package datatypes

import (
	"time"

	"github.com/edgedb/edgedb-go"
)

type RequestErrorInterface interface {
	Error() string
	StatusCode() int
}

type RequestError struct {
	statusCode int
	message    string
}

func (e *RequestError) Error() string {
	return e.message
}

func (e *RequestError) StatusCode() int {
	return e.statusCode
}

func NewRequestError(statusCode int, message string) error {
	return &RequestError{
		statusCode: statusCode,
		message:    message,
	}
}

type Token struct {
	ID        edgedb.UUID `json:"id" edgedb:"id"`
	Variant   string      `json:"variant" edgedb:"variant"`
	Value     string      `json:"value" edgedb:"value"`
	Scope     []string    `json:"scope" edgedb:"scope"`
	Account   Account     `json:"account" edgedb:"account"`
	Revoked   bool        `json:"revoked" edgedb:"revoked"`
	ExpiresAt time.Time   `json:"expires_at" edgedb:"expires_at"`
}

type Password struct {
	Id                edgedb.UUID `edgedb:"id"`
	Account           Account     `edgedb:"account"`
	Email             string      `edgedb:"email"`
	Password          string      `edgedb:"password"`
	LastUsed          time.Time   `edgedb:"last_used"`
	FailedAttempts    int16       `edgedb:"failed_attempts"`
	LastFailedAttempt time.Time   `edgedb:"last_failed_attempt"`
}

type Account struct {
	Id                edgedb.UUID        `edgedb:"id"`
	Username          string             `edgedb:"username"`
	Email             string             `edgedb:"email"`
	AvatarURI         string             `edgedb:"avatar_uri"`
	Status            string             `edgedb:"status"`
	OtpSecret         edgedb.OptionalStr `edgedb:"otp_secret"`
	OtpState          string             `edgedb:"otp_state"`
	StatusDescription string             `edgedb:"status_description"`
	StatusChanged     time.Time          `edgedb:"status_changed"`
	CreatedAt         time.Time          `edgedb:"created_at"`
}
