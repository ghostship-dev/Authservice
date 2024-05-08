package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/queries"
	"github.com/ghostship-dev/authservice/core/responses"
	"github.com/ghostship-dev/authservice/core/utility"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.LoginRequestData

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return responses.NewErrorResponse(w, http.StatusBadRequest, "request was malformed or invalid")
	}

	if err := reqData.Validate(); err != nil {
		return responses.NewErrorResponse(w, http.StatusBadRequest, "request was invalid")
	}

	password, err := queries.GetPasswordByEmail(reqData.Email)
	if err != nil {
		var edbErr edgedb.Error
		if errors.As(err, &edbErr) && edbErr.Category(edgedb.NoDataError) {
			return responses.NewErrorResponse(w, http.StatusBadRequest, "account not found or wrong password")
		}
		return responses.NewErrorResponse(w, http.StatusInternalServerError, "internal server error")
	}

	if password.FailedAttempts >= 3 {
		return responses.NewErrorResponse(w, http.StatusBadRequest, "account has been suspended, due to too many failed login attempts")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password.Password), []byte(reqData.Password)); err != nil {
		err = queries.IncrementFailedPasswordLoginAttempts(reqData.Email)
		if err != nil {
			fmt.Println(err)
		}
		return responses.NewErrorResponse(w, http.StatusBadRequest, "account not found or wrong password")
	}

	//TODO: Add Two Factor Authentication (Google Authenticator)

	if password.FailedAttempts > 0 {
		err = queries.ResetFailedPasswordLoginAttempts(reqData.Email)
		if err != nil {
			return responses.NewErrorResponse(w, http.StatusInternalServerError, "internal server error")
		}
	}

	accessToken, err := utility.NewAccessToken(password.Account.Id, make([]string, 0))
	if err != nil {
		return responses.NewErrorResponse(w, http.StatusInternalServerError, "internal server error")
	}

	refreshToken, err := utility.NewRefreshToken(password.Account.Id, make([]string, 0))
	if err != nil {
		return responses.NewErrorResponse(w, http.StatusInternalServerError, "internal server error")
	}

	if err = queries.AddNewTokenPair(password.Account.Id, accessToken.Value, refreshToken.Value, time.Now().Add(time.Hour*2), time.Now().Add(time.Hour*24*30), accessToken.Scope); err != nil {
		return responses.NewErrorResponse(w, http.StatusInternalServerError, "internal server error")
	}

	return responses.SendLoginSuccessResponse(accessToken, refreshToken, w)
}
