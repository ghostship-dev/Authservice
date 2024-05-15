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

func OAuth2ApplicationNameInUseErrorResponse() error {
	response := LoginErrorResponse{
		Error:       true,
		Message:     "oauth2_application_name_in_use",
		Description: "OAuth2 application name is already in use",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(http.StatusOK, string(jsonResponse))
}

func OAuth2ApplicationTypeParameterNameMismatchErrorResponse() error {
	response := LoginErrorResponse{
		Error:       true,
		Message:     "oauth2_application_type_parameter_name_mismatch",
		Description: "OAuth2 application type parameter name mismatch",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(http.StatusOK, string(jsonResponse))
}
