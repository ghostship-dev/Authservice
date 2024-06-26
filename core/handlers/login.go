package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/database"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/responses"
	"github.com/ghostship-dev/authservice/core/utility"
	"github.com/xlzd/gotp"
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

	password, err := database.Connection.Queries.GetPasswordByEmail(reqData.Email)
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
		err = database.Connection.Queries.IncrementFailedPasswordLoginAttempts(reqData.Email)
		if err != nil {
			fmt.Println(err)
		}
		return responses.AccountNotFoundResponse()
	}

	if password.Account.OtpState == "enabled" {
		if len(reqData.OTP) < 6 {
			return responses.TwoFactorAuthenticationRequiredResponse()
		} else {
			totpSecret, isSecretSet := password.Account.OtpSecret.Get()
			if !isSecretSet {
				return responses.InvalidTotpStateErrorResponse()
			}

			totp := gotp.NewDefaultTOTP(totpSecret)
			if !totp.Verify(reqData.OTP, time.Now().Unix()) {
				return responses.UnauthorizedErrorResponse("invalid otp")
			}
		}
	}

	if password.FailedAttempts > 0 {
		err = database.Connection.Queries.ResetFailedPasswordLoginAttempts(reqData.Email)
		if err != nil {
			return responses.InternalServerErrorResponse()
		}
	}

	accessTokenExpiresAt := time.Now().Add(time.Hour * 1)
	refreshTokenExpiresAt := time.Now().Add(time.Hour * 24 * 7)

	grantedScope := []string{"*"}

	accessToken, err := utility.NewAccessToken(password.Account, accessTokenExpiresAt, grantedScope)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	refreshToken, err := utility.NewRefreshToken(password.Account, refreshTokenExpiresAt, grantedScope)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	if err = database.Connection.Queries.AddNewTokenPair(password.Account.Id, accessToken.Value, refreshToken.Value, accessTokenExpiresAt, refreshTokenExpiresAt, accessToken.Scope); err != nil {
		return responses.InternalServerErrorResponse()
	}

	return responses.SendLoginSuccessResponse(accessToken, refreshToken, w)
}
