package responses

import (
	"encoding/json"
	"net/http"

	"github.com/ghostship-dev/authservice/core/datatypes"
)

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

func ValidationErrorResponse(errors map[string]string) error {
	jsonErrors, _ := json.Marshal(errors)
	return datatypes.NewRequestError(http.StatusBadRequest, string(jsonErrors))
}

func InternalServerErrorResponse() error {
	return datatypes.NewRequestError(http.StatusInternalServerError, "internal server error")
}
