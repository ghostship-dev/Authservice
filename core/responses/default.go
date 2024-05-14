package responses

import (
	"encoding/json"
	"net/http"

	"github.com/ghostship-dev/authservice/core/datatypes"
)

type GenericResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func NewJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func NewErrorResponse(w http.ResponseWriter, statusCode int, message string) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(map[string]string{"error": message, "success": "false"})
}

func NewTextResponse(w http.ResponseWriter, statusCode int, message string) error {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)

	if _, err := w.Write([]byte(message)); err != nil {
		return InternalServerErrorResponse()
	}

	return nil
}

func BadRequestResponse() error {
	return datatypes.NewRequestError(http.StatusBadRequest, "bad request (request was malformed)")
}

type ValidationErrorResponseType struct {
	Error  bool              `json:"error"`
	Fields map[string]string `json:"fields"`
}

func ValidationErrorResponse(errors map[string]string) error {
	response := ValidationErrorResponseType{
		Error:  true,
		Fields: errors,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(http.StatusBadRequest, string(jsonResponse))
}

type UnauthorizedErrorResponseType struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func UnauthorizedErrorResponse(message string) error {
	response := UnauthorizedErrorResponseType{
		Error:   true,
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(http.StatusUnauthorized, string(jsonResponse))
}

func InternalServerErrorResponse() error {
	return datatypes.NewRequestError(http.StatusInternalServerError, "internal server error")
}

func SendNewOKResponse(w http.ResponseWriter) error {
	err := NewJSONResponse(w, http.StatusOK, GenericResponse{
		Error:   false,
		Message: "action successful",
	})
	if err != nil {
		return InternalServerErrorResponse()
	}
	return nil
}
