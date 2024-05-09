package responses

import (
	"encoding/json"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"net/http"
)

type RegisterSuccessResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func SendRegisterSuccessResponse(w http.ResponseWriter) error {
	err := NewJSONResponse(w, http.StatusOK, RegisterSuccessResponse{
		Error:   false,
		Message: "account successfully created",
	})
	if err != nil {
		return InternalServerErrorResponse()
	}
	return nil
}

func EmailInUseErrorResponse() error {
	response := LoginErrorResponse{
		Error:       true,
		Message:     "email_in_use",
		Description: "Email already in use",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(http.StatusOK, string(jsonResponse))
}

func InvalidEmailErrorResponse() error {
	response := LoginErrorResponse{
		Error:       true,
		Message:     "invalid_email",
		Description: "Provided email is invalid",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(http.StatusOK, string(jsonResponse))
}
