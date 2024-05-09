package handlers

import (
	"encoding/json"
	"errors"
	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/queries"
	"github.com/ghostship-dev/authservice/core/responses"
	"io"
	"net/http"
	"strings"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.RegisterRequestData

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return responses.NewErrorResponse(w, http.StatusBadRequest, "request was malformed or invalid")
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if err := reqData.Validate(); err != nil {
		return responses.NewErrorResponse(w, http.StatusBadRequest, "request was malformed or invalid")
	}

	account, err := queries.CreateAccount(reqData.Email, reqData.Username, reqData.Password)
	if err != nil {
		var edbErr edgedb.Error
		if errors.As(err, &edbErr) && edbErr.Category(edgedb.ConstraintViolationError) {
			if strings.Contains(edbErr.Error(), "violates exclusivity constraint") {
				return responses.NewErrorResponse(w, http.StatusBadRequest, "email already taken")
			}
			if strings.Contains(edbErr.Error(), "invalid email") {
				return responses.NewErrorResponse(w, http.StatusBadRequest, "email invalid")
			}
		}
		return responses.NewErrorResponse(w, http.StatusInternalServerError, "internal server error")
	}

	_ = account

	return responses.NewTextResponse(w, http.StatusOK, "ok for now")
}
