package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/responses"
)

func Interosppect(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.IntrospectRequestData

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return responses.NewErrorResponse(w, http.StatusBadRequest, "request was malformed")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if err := reqData.Validate(); err != nil {
		return responses.NewErrorResponse(w, http.StatusBadRequest, "request was invalid")
	}

	return nil
}
