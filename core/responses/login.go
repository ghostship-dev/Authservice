package responses

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ghostship-dev/authservice/core/datatypes"
)

type LoginSuccessResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func SendLoginSuccessResponse(accessToken, refreshToken datatypes.Token, w http.ResponseWriter) error {
	err := NewJSONResponse(w, http.StatusOK, LoginSuccessResponse{
		AccessToken:  accessToken.Value,
		TokenType:    "Bearer",
		ExpiresAt:    accessToken.ExpiresAt,
		RefreshToken: refreshToken.Value,
	})
	if err != nil {
		return InternalServerErrorResponse()
	}
	return nil
}

type LoginErrorResponse struct {
	Error       bool   `json:"error"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func AccountNotFoundResponse() error {
	response := LoginErrorResponse{
		Error:       true,
		Message:     "account_not_found",
		Description: "Account not found",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(http.StatusOK, string(jsonResponse))
}

func ToManyFailedAttemptsResponse() error {
	response := LoginErrorResponse{
		Error:       true,
		Message:     "to_many_failed_attempts",
		Description: "Account has been suspended due to too many failed login attempts",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(http.StatusOK, string(jsonResponse))
}
