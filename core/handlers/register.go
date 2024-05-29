package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/database"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/responses"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.RegisterRequestData

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return responses.BadRequestResponse()
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if validationErrors := reqData.Validate(); len(validationErrors) > 0 {
		return responses.ValidationErrorResponse(validationErrors)
	}

	account, err := database.Connection.Queries.CreateAccount(reqData.Email, reqData.Username, reqData.Password)
	if err != nil {
		var edbErr edgedb.Error
		if errors.As(err, &edbErr) && edbErr.Category(edgedb.ConstraintViolationError) {
			if strings.Contains(edbErr.Error(), "violates exclusivity constraint") {
				return responses.EmailInUseErrorResponse()
			}
			if strings.Contains(edbErr.Error(), "invalid email") {
				return responses.InvalidEmailErrorResponse()
			}
		}
		return responses.InternalServerErrorResponse()
	}

	_ = account

	return responses.SendRegisterSuccessResponse(w)
}
