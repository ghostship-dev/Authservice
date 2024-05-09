package responses

import (
	"github.com/ghostship-dev/authservice/core/datatypes"
	"net/http"
	"time"
)

type LoginSuccessResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func SendLoginSuccessResponse(accessToken, refreshToken datatypes.Token, w http.ResponseWriter) error {
	return NewJSONResponse(w, http.StatusOK, LoginSuccessResponse{
		AccessToken:  accessToken.Value,
		TokenType:    "Bearer",
		ExpiresAt:    accessToken.ExpiresAt,
		RefreshToken: refreshToken.Value,
	})
}
