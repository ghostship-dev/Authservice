package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/queries"
	"github.com/ghostship-dev/authservice/core/responses"
	"github.com/ghostship-dev/authservice/core/utility"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.LoginRequestData

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return responses.BadRequestResponse()
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if validationErrors := reqData.Validate(); len(validationErrors) > 0 {
		return responses.ValidationErrorResponse(validationErrors)
	}

	password, err := queries.GetPasswordByEmail(reqData.Email)
	if err != nil {
		var edbErr edgedb.Error
		if errors.As(err, &edbErr) && edbErr.Category(edgedb.NoDataError) {
			return responses.AccountNotFoundResponse()
		}
		return responses.InternalServerErrorResponse()
	}

	if password.FailedAttempts >= 3 {
		return responses.ToManyFailedAttemptsResponse()
	}

	if err = bcrypt.CompareHashAndPassword([]byte(password.Password), []byte(reqData.Password)); err != nil {
		err = queries.IncrementFailedPasswordLoginAttempts(reqData.Email)
		if err != nil {
			fmt.Println(err)
		}
		return responses.AccountNotFoundResponse()
	}

	//TODO: Add Two Factor Authentication (Google Authenticator)

	if password.FailedAttempts > 0 {
		err = queries.ResetFailedPasswordLoginAttempts(reqData.Email)
		if err != nil {
			return responses.InternalServerErrorResponse()
		}
	}

	accessTokenExpiresAt := time.Now().Add(time.Hour * 1)
	refreshTokenExpiresAt := time.Now().Add(time.Hour * 24 * 7)

	grantedScope := []string{"account_*"}

	accessToken, err := utility.NewAccessToken(password.Account, accessTokenExpiresAt, grantedScope)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	refreshToken, err := utility.NewRefreshToken(password.Account, refreshTokenExpiresAt, grantedScope)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	if err = queries.AddNewTokenPair(password.Account.Id, accessToken.Value, refreshToken.Value, accessTokenExpiresAt, refreshTokenExpiresAt, accessToken.Scope); err != nil {
		return responses.InternalServerErrorResponse()
	}

	return responses.SendLoginSuccessResponse(accessToken, refreshToken, w)
}
