package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/database"
	"github.com/ghostship-dev/authservice/core/scopes"
	"github.com/ghostship-dev/authservice/core/utility"
	"github.com/golang-jwt/jwt/v5"
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

	dbToken, err := database.Connection.Queries.GetToken(bearerToken)
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

	if err = database.Connection.Queries.CreateNewOAuthClientApplication(oauthApplication); err != nil {
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

	response := responses.GenericDataResponse{
		Error: false,
		Data: struct {
			ClientID string `json:"client_id"`
		}{
			ClientID: oauthApplication.ClientID,
		},
	}

	return responses.NewJSONResponse(w, http.StatusOK, response)
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

	dbToken, err := database.Connection.Queries.GetToken(bearerToken)
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

	if err = database.Connection.Queries.UpdateOAuth2ClientApplicationKeyValue(reqData); err != nil {
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

func DeleteOAuthClientApplication(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.DeleteOAuth2ClientRequest
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

	dbToken, err := database.Connection.Queries.GetToken(bearerToken)
	if err != nil {
		return responses.UnauthorizedErrorResponse("invalid bearer token")
	}

	if dbToken.Revoked {
		return responses.UnauthorizedErrorResponse("token is revoked")
	}

	// TODO: Replace with a permission decision point in the future
	if !slices.Contains(dbToken.Scope, "oauth2_delete") && !slices.Contains(dbToken.Scope, "*") {
		return responses.UnauthorizedErrorResponse("missing required permission")
	}

	if err = database.Connection.Queries.DeleteOAuth2ClientApplication(reqData.ClientID); err != nil {
		return responses.InternalServerErrorResponse()
	}

	return responses.SendNewOKResponse(w)
}

func IntrospectOAuthToken(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.IntrospectOAuth2TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return responses.BadRequestResponse()
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if validationErrors := reqData.Validate(); len(validationErrors) > 0 {
		return responses.ValidationErrorResponse(validationErrors)
	}

	token, err := jwt.Parse(reqData.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
		return responses.UnauthorizedErrorResponse("invalid token")
	}

	if reqData.CheckIfTokenIsRevoked {
		dbToken, err := database.Connection.Queries.GetToken(reqData.Token)
		tokenVariant := token.Claims.(jwt.MapClaims)["variant"].(string)
		if err != nil {
			return responses.UnauthorizedErrorResponse("invalid token")
		}
		if dbToken.Variant != tokenVariant {
			return responses.UnauthorizedErrorResponse("token type mismatch")
		}
		if dbToken.Revoked {
			return responses.UnauthorizedErrorResponse("token is revoked")
		}
	}

	return responses.SendNewOKResponseMessage(w, "token is valid")
}

func AuthorizeOAuthApplication(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return responses.BadRequestResponse()
	}

	reqData := datatypes.AuthorizeOAuth2ClientRequest{
		ClientID:     r.Form.Get("client_id"),
		UserID:       r.Form.Get("user_id"),
		RedirectURI:  r.Form.Get("redirect_uri"),
		ResponseType: r.Form.Get("response_type"),
		Scope:        r.Form.Get("scope"),
		State:        r.Form.Get("state"),
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if validationErrors := reqData.Validate(); len(validationErrors) > 0 {
		return responses.ValidationErrorResponse(validationErrors)
	}

	userID, err := edgedb.ParseUUID(reqData.UserID)
	if err != nil {
		return responses.BadRequestResponse()
	}

	oauth2Application, account, err := database.Connection.Queries.GetOAuth2ClientApplicationAndUserAccount(reqData.ClientID, userID)
	if err != nil {
		return responses.OAuth2ApplicationNotFoundResponse()
	}

	if len(oauth2Application.ClientID) < 1 {
		return responses.OAuth2ApplicationNotFoundResponse()
	}

	if len(account.Id.String()) < 1 {
		return responses.OAuth2UserNotFoundResponse()
	}

	if !slices.Contains(oauth2Application.RedirectURIs, reqData.RedirectURI) {
		return responses.OAuth2RedirectURIDoesNotMatch()
	}

	scopeSlice := strings.Split(strings.ReplaceAll(reqData.Scope, " ", ""), ",")

	if len(scopeSlice) < 1 {
		return responses.OAuth2ScopeIsRequired()
	}

	if !scopes.AllScopesAllowed(scopeSlice) {
		return responses.OAuth2InvalidScope(scopes.GetForbiddenScopes(scopeSlice))
	}

	stateToken, err := gonanoid.New(50)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	var authCode = datatypes.OAuthAuthorizationCode{
		Code:           stateToken,
		Consented:      false,
		ExpiresAt:      time.Now().Add(10 * time.Minute),
		GrantedScope:   make([]string, 0),
		RequestedScope: scopeSlice,
		Account:        account,
		Application:    oauth2Application,
		RedirectURI:    reqData.RedirectURI,
	}

	if err = database.Connection.Queries.CreateNewOAuth2AuthorizationCode(authCode); err != nil {
		return responses.InternalServerErrorResponse()
	}

	return responses.ReturnRedirectResponseToConsentPage(w, r, authCode)
}

func OAuthTokenEndpoint(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.OAuthTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return responses.BadRequestResponse()
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if validationErrors := reqData.Validate(); len(validationErrors) > 0 {
		return responses.ValidationErrorResponse(validationErrors)
	}

	clientSecret, err := utility.GetClientSecretFromHeader(&r.Header)
	if err != nil || len(clientSecret) < 1 {
		return responses.UnauthorizedErrorResponse("missing client secret")
	}
	clientID, err := utility.GetClientIDFromHeader(&r.Header)
	if err != nil || len(clientID) < 1 {
		return responses.UnauthorizedErrorResponse("missing client id")
	}

	if reqData.GrantType == "authorization_code" {
		return handleAuthorizationCodeGrantType(w, reqData, clientID, clientSecret)
	}

	if reqData.GrantType == "refresh_token" {
		return handleRefreshTokenGrantType(w, reqData)
	}

	return responses.BadRequestResponse()
}

func handleAuthorizationCodeGrantType(w http.ResponseWriter, reqData datatypes.OAuthTokenRequest, clientID, clientSecret string) error {
	authCode, err := database.Connection.Queries.GetOAuth2AuthorizationCode(reqData.Code)
	if err != nil {
		return responses.UnauthorizedErrorResponse("invalid authorization code")
	}

	if authCode.Application.ClientID != clientID {
		return responses.UnauthorizedErrorResponse("invalid client id")
	}

	if authCode.Application.ClientSecret != clientSecret {
		return responses.UnauthorizedErrorResponse("invalid client secret")
	}

	if authCode.ExpiresAt.Before(time.Now()) {
		return responses.UnauthorizedErrorResponse("authorization code expired")
	}

	if !authCode.Consented {
		return responses.UnauthorizedErrorResponse("authorization code not consented")
	}

	accessTokenExpiresAt := time.Now().Add(time.Hour * 1)
	refreshTokenExpiresAt := time.Now().Add(time.Hour * 24 * 7)

	accessToken, err := utility.NewAccessToken(authCode.Account, accessTokenExpiresAt, authCode.GrantedScope)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	refreshToken, err := utility.NewRefreshToken(authCode.Account, refreshTokenExpiresAt, authCode.GrantedScope)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	if err = database.Connection.Queries.AddNewTokenPair(authCode.Account.Id, accessToken.Value, refreshToken.Value, accessTokenExpiresAt, refreshTokenExpiresAt, accessToken.Scope); err != nil {
		return responses.InternalServerErrorResponse()
	}

	if err = database.Connection.Queries.DeleteOAuth2AuthorizationCode(authCode.Code); err != nil {
		return responses.InternalServerErrorResponse()
	}

	return responses.SendTokenExchangeSuccessResponse(accessToken, refreshToken, w)
}

func handleRefreshTokenGrantType(w http.ResponseWriter, reqData datatypes.OAuthTokenRequest) error {
	refreshToken, err := database.Connection.Queries.GetRefreshToken(reqData.RefreshToken)
	if err != nil {
		fmt.Println(err)
		return responses.UnauthorizedErrorResponse("invalid refresh token")
	}

	if refreshToken.Variant != "refresh_token" {
		fmt.Println(refreshToken.Variant)
		return responses.UnauthorizedErrorResponse("invalid refresh token")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return responses.UnauthorizedErrorResponse("refresh token expired")
	}

	accessTokenExpiresAt := time.Now().Add(time.Hour * 1)
	refreshTokenExpiresAt := time.Now().Add(time.Hour * 24 * 7)

	accessToken, err := utility.NewAccessToken(refreshToken.Account, accessTokenExpiresAt, refreshToken.Scope)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	newRefreshToken, err := utility.NewRefreshToken(refreshToken.Account, refreshTokenExpiresAt, refreshToken.Scope)
	if err != nil {
		return responses.InternalServerErrorResponse()
	}

	if err = database.Connection.Queries.AddNewTokenPair(refreshToken.Account.Id, accessToken.Value, newRefreshToken.Value, accessTokenExpiresAt, refreshTokenExpiresAt, accessToken.Scope); err != nil {
		return responses.InternalServerErrorResponse()
	}

	if err = database.Connection.Queries.DeleteRefreshToken(refreshToken.ID); err != nil {
		return responses.InternalServerErrorResponse()
	}

	return responses.SendTokenExchangeSuccessResponse(accessToken, newRefreshToken, w)
}

func RevokeOAuthToken(w http.ResponseWriter, r *http.Request) error {
	var reqData datatypes.OAuthRevokeTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return responses.BadRequestResponse()
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	if validationErrors := reqData.Validate(); len(validationErrors) > 0 {
		return responses.ValidationErrorResponse(validationErrors)
	}

	if err := database.Connection.Queries.DeleteTokensByValue([]string{reqData.Token}); err != nil {
		fmt.Println(err)
		return responses.InternalServerErrorResponse()
	}

	return responses.SendNewOKResponseMessage(w, "token revoked successfully")
}
