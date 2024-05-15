package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/queries"
	"github.com/ghostship-dev/authservice/core/utility"
	gonanoid "github.com/matoous/go-nanoid/v2"

	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/responses"
)

func NewOAuthApplication(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.NewOAuthClientRequest

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
		return responses.UnauthorizedErrorResponse("invalid bearer token")
	}

	if dbToken.Revoked {
		return responses.UnauthorizedErrorResponse("token is revoked")
	}

	// TODO: Replace with a permission decision point in the future
	if !slices.Contains(dbToken.Scope, "oauth2_write") && !slices.Contains(dbToken.Scope, "*") {
		return responses.UnauthorizedErrorResponse("missing required permission")
	}

	// TODO: create new oauth2 client application

	clientId, err := gonanoid.New(30)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}
	clientSecret, err := gonanoid.New(30)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	oauthApplication := datatypes.OAuthClient{
		ClientID:               clientId,
		ClientSecret:           clientSecret,
		ClientName:             reqData.ClientName,
		ClientType:             reqData.ClientType,
		RedirectURIs:           reqData.RedirectUris,
		GrantTypes:             reqData.GrantTypes,
		Scope:                  reqData.Scope,
		ClientOwner:            dbToken.Account,
		ClientDescription:      edgedb.NewOptionalStr(reqData.ClientDescription),
		ClientHomepageUrl:      edgedb.NewOptionalStr(reqData.ClientHomepageUrl),
		ClientLogoUrl:          edgedb.NewOptionalStr(reqData.ClientLogoUrl),
		ClientTosUrl:           edgedb.NewOptionalStr(reqData.ClientTosUrl),
		ClientPrivacyUrl:       edgedb.NewOptionalStr(reqData.ClientPrivacyUrl),
		ClientRegistrationDate: time.Now(),
		ClientStatus:           "active",
		ClientRateLimits:       []byte(""),
	}

	if err = queries.CreateNewOAuthClientApplication(oauthApplication); err != nil {
		var edbErr edgedb.Error
		if errors.As(err, &edbErr) && edbErr.Category(edgedb.ConstraintViolationError) {
			if strings.Contains(edbErr.Error(), "violates exclusivity constraint") {
				return responses.OAuth2ApplicationNameInUseErrorResponse()
			} else {
				return responses.BadRequestResponse()
			}
		}
		return responses.InternalServerErrorResponse()
	}

	return nil
}

func UpdateOAuthApplicationKeyValue(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.UpdateOAuth2ClientKeyValueRequest
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
		return responses.UnauthorizedErrorResponse("invalid bearer token")
	}

	if dbToken.Revoked {
		return responses.UnauthorizedErrorResponse("token is revoked")
	}

	// TODO: Replace with a permission decision point in the future
	if !slices.Contains(dbToken.Scope, "oauth2_write") && !slices.Contains(dbToken.Scope, "*") {
		return responses.UnauthorizedErrorResponse("missing required permission")
	}

	if err = queries.UpdateOAuth2ClientApplicationKeyValue(reqData); err != nil {
		var edbErr edgedb.Error
		if errors.As(err, &edbErr) && edbErr.Category(edgedb.ConstraintViolationError) {
			if strings.Contains(edbErr.Error(), "violates exclusivity constraint") {
				return responses.OAuth2ApplicationNameInUseErrorResponse()
			} else {
				return responses.BadRequestResponse()
			}
		} else if edbErr.Category(edgedb.ParameterTypeMismatchError) {
			return responses.OAuth2ApplicationTypeParameterNameMismatchErrorResponse()
		}
		return responses.InternalServerErrorResponse()
	}

	return responses.SendNewOKResponse(w)
}
