package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/ghostship-dev/authservice/core/queries"
	"github.com/ghostship-dev/authservice/core/utility"
	"github.com/xlzd/gotp"

	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/responses"
)

func AccountOTP(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.OTPRequest

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return responses.BadRequestResponse()
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if validationErrors := reqData.Validate(); len(validationErrors) > 0 {
		return responses.ValidationErrorResponse(validationErrors)
	}

	bearerToken, err := utility.GetBearerTokenFromHeader(&r.Header)
	if err != nil {
		return responses.UnauthorizedErrorResponse("missing bearer token")
	}

	dbToken, err := queries.GetToken(bearerToken)
	if err != nil {
		fmt.Println(err)
		return responses.UnauthorizedErrorResponse("invalid bearer token")
	}

	if dbToken.Revoked {
		return responses.UnauthorizedErrorResponse("token is revoked")
	}

	// TODO: Replace with a permission decision point in the future
	if !slices.Contains(dbToken.Scope, "account_otp_write") && !slices.Contains(dbToken.Scope, "account_*") {
		return responses.UnauthorizedErrorResponse("missing required permission")
	}

	// Enable Section
	if reqData.Action == "enable" && dbToken.Account.OtpState == "disabled" {
		totpSecret := gotp.RandomSecret(16)
		totp := gotp.NewDefaultTOTP(totpSecret)
		qrUri := totp.ProvisioningUri(dbToken.Account.Username, "Ghostship")

		if err = queries.SetOTPSecret(dbToken.Account.Id, totpSecret); err != nil {
			fmt.Println(err)
			return responses.InternalServerErrorResponse()
		}

		return responses.SendActivateTotpSuccessResponse(totpSecret, qrUri, w)
	}

	secret, secretFound := dbToken.Account.OtpSecret.Get()
	if !secretFound {
		return responses.UnauthorizedErrorResponse("invalid otp secret")
	}

	// Disable Section
	if reqData.Action == "disable" && dbToken.Account.OtpState == "enabled" {
		totp := gotp.NewDefaultTOTP(secret)
		if !totp.Verify(reqData.OTP, time.Now().Unix()) {
			return responses.UnauthorizedErrorResponse("invalid otp")
		}

		if err = queries.ResetOTP(dbToken.Account.Id); err != nil {
			return responses.InternalServerErrorResponse()
		}

		return responses.SendDisableTotpSuccessResponse(w)
	}

	// Verify Section
	if reqData.Action == "verify" && dbToken.Account.OtpState == "verifying" {
		totp := gotp.NewDefaultTOTP(secret)
		if !totp.Verify(reqData.OTP, time.Now().Unix()) {
			return responses.UnauthorizedErrorResponse("invalid otp")
		}

		if err = queries.SetOTPState(dbToken.Account.Id, "enabled"); err != nil {
			return responses.InternalServerErrorResponse()
		}

		return responses.SendEnableTotpSuccessResponse(w)
	}

	return responses.InvalidTotpStateErrorResponse()
}
