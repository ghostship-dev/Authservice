package responses

import (
	"encoding/json"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"net/http"
)

type ActivateTotpSuccessResponse struct {
	Error  bool   `json:"error"`
	Secret string `json:"secret"`
	QrUri  string `json:"qr_uri"`
}

func SendActivateTotpSuccessResponse(secret, uri string, w http.ResponseWriter) error {
	err := NewJSONResponse(w, http.StatusOK, ActivateTotpSuccessResponse{
		Error:  false,
		Secret: secret,
		QrUri:  uri,
	})
	if err != nil {
		return InternalServerErrorResponse()
	}
	return nil
}

func SendEnableTotpSuccessResponse(w http.ResponseWriter) error {
	err := NewJSONResponse(w, http.StatusOK, GenericResponse{
		Error:   false,
		Message: "totp enabled",
	})
	if err != nil {
		return InternalServerErrorResponse()
	}
	return nil
}

func SendDisableTotpSuccessResponse(w http.ResponseWriter) error {
	err := NewJSONResponse(w, http.StatusOK, GenericResponse{
		Error:   false,
		Message: "totp disabled",
	})
	if err != nil {
		return InternalServerErrorResponse()
	}
	return nil
}

type ActivateTotpErrorResponse struct {
	Error       bool   `json:"error"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func InvalidTotpStateErrorResponse() error {
	response := ActivateTotpErrorResponse{
		Error:       true,
		Message:     "invalid_totp_state",
		Description: "totp state is invalid",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(http.StatusOK, string(jsonResponse))
}
