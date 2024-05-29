package queries

import (
	"context"
	"errors"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/database"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/responses"
	"golang.org/x/crypto/bcrypt"
)

func GetPasswordByEmail(email string) (datatypes.Password, error) {
	var password datatypes.Password
	query := "SELECT Password{password, failed_attempts, account: { id, otp_secret, otp_state }} filter .email = <str>$0 LIMIT 1"
	err := database.Client.QuerySingle(database.Context, query, &password, email)
	return password, err
}

func CreateAccount(email, username, password string) (datatypes.Account, error) {
	var account datatypes.Account
	err := database.Client.Tx(database.Context, func(ctx context.Context, tx *edgedb.Tx) error {
		accountCreationQuery := "INSERT Account { username := <str>$0, email := <str>$1 }"
		passwordCreationQuery := "INSERT Password { account := <Account>$0, email := <str>$1, password := <str>$2 }"

		err := tx.QuerySingle(ctx, accountCreationQuery, &account, username, email)
		if err != nil {
			return err
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		err = tx.Execute(ctx, passwordCreationQuery, account.Id, email, string(hashedPassword))
		if err != nil {
			return err
		}

		return nil
	})
	return account, err
}

func GetAccountById(id string) (datatypes.Account, error) {
	var account datatypes.Account
	query := "SELECT Account filter .id = <uuid>$0 LIMIT 1"
	accountId, err := edgedb.ParseUUID(id)
	if err != nil {
		return account, err
	}
	err = database.Client.QuerySingle(database.Context, query, &account, accountId)
	return account, err
}

func IncrementFailedPasswordLoginAttempts(email string) error {
	query := "UPDATE Password filter .email = <str>$0 set { failed_attempts := .failed_attempts +1, last_failed_attempt := <datetime>$1 }"
	return database.Client.Execute(database.Context, query, email, time.Now())
}

func ResetFailedPasswordLoginAttempts(email string) error {
	query := "UPDATE Password filter .email = <str>$0 set { failed_attempts := 0, last_failed_attempt := {} }"
	return database.Client.Execute(database.Context, query, email)
}

func AddNewToken(accountId edgedb.UUID, value, variant string, expiresAt time.Time, scope []string) error {
	query := "INSERT Token { account := <Account>$0, variant := <str>$1, scope := <array<str>>$2, value := <str>$3, revoked := <bool>$4, expires_at := <datetime>$5 }"
	return database.Client.Execute(database.Context, query, accountId, variant, scope, value, false, expiresAt)
}

func AddNewTokenPair(accountId edgedb.UUID, accessTokenValue, refreshTokenValue string, accessTokenExpiresAt, refreshTokenExpiresAt time.Time, scope []string) error {
	query := "INSERT Token { account := <Account>$0, variant := <str>$1, scope := <array<str>>$2, value := <str>$3, revoked := <bool>$4, expires_at := <datetime>$5 }; INSERT Token { account := <Account>$6, variant := <str>$7, scope := <array<str>>$8, value := <str>$9, revoked := <bool>$10, expires_at := <datetime>$11 }"
	return database.Client.Execute(database.Context, query, accountId, "access_token", scope, accessTokenValue, false, accessTokenExpiresAt, accountId, "refresh_token", scope, refreshTokenValue, false, refreshTokenExpiresAt)
}

func GetToken(tokenValue string) (datatypes.Token, error) {
	var token datatypes.Token
	query := "SELECT Token { value, scope, revoked, variant, expires_at, account: { id, username, otp_secret, otp_state } } filter .value = <str>$0 LIMIT 1"
	return token, database.Client.QuerySingle(database.Context, query, &token, tokenValue)
}

func ResetOTP(accountId edgedb.UUID) error {
	query := "UPDATE Account filter .id = <uuid>$0 set { otp_secret := <str>{}, otp_state := <str>'disabled' }"
	return database.Client.Execute(database.Context, query, accountId)
}

func SetOTPSecret(accountId edgedb.UUID, otpSecret string) error {
	query := "UPDATE Account filter .id = <uuid>$0 set { otp_secret := <str>$1, otp_state := <str>$2 }"
	return database.Client.Execute(database.Context, query, accountId, otpSecret, "verifying")
}

func SetOTPState(accountId edgedb.UUID, otpState string) error {
	if otpState != "disabled" && otpState != "enabled" && otpState != "verifying" {
		return errors.New("allowed otp state: disabled, enabled or verifying")
	}
	query := "UPDATE Account filter .id = <uuid>$0 set { otp_state := <str>$1 }"
	return database.Client.Execute(database.Context, query, accountId, otpState)
}

func CreateNewOAuthClientApplication(oauthClient datatypes.OAuthClient) error {
	query := `
		INSERT OAuthApplication {
			client_id := <str>$0,
			client_secret := <str>$1,
			client_name := <str>$2,
			client_type := <str>$3,
			redirect_uris := <array<str>>$4,
			grant_types := <array<str>>$5,
			scope := <array<str>>$6,
			client_owner := <Account>$7,
			client_description := <str>$8,
			client_homepage_url := <str>$9,
			client_logo_url := <str>$10,
			client_tos_url := <str>$11,
			client_privacy_url := <str>$12,
			client_registration_date := <datetime>$13,
			client_status := <str>$14,
		}
	`

	return database.Client.Execute(database.Context, query,
		oauthClient.ClientID,
		oauthClient.ClientSecret,
		oauthClient.ClientName,
		oauthClient.ClientType,
		oauthClient.RedirectURIs,
		oauthClient.GrantTypes,
		oauthClient.Scope,
		oauthClient.ClientOwner.Id,
		oauthClient.ClientDescription,
		oauthClient.ClientHomepageUrl,
		oauthClient.ClientLogoUrl,
		oauthClient.ClientTosUrl,
		oauthClient.ClientPrivacyUrl,
		oauthClient.ClientRegistrationDate,
		oauthClient.ClientStatus,
	)
}

func UpdateOAuth2ClientApplicationKeyValue(updateRequestData datatypes.UpdateOAuth2ClientKeyValueRequest) error {
	keyType, err := updateRequestData.GetKeyType()
	if err != nil {
		return responses.BadRequestResponse()
	}
	query := "UPDATE OAuthApplication filter .client_id = <str>$0 set { " + updateRequestData.Key + " := " + keyType + "'" + updateRequestData.Value + "' }"
	return database.Client.Execute(database.Context, query, updateRequestData.ClientID)
}

func DeleteOAuth2ClientApplication(clientId string) error {
	query := "DELETE OAuthApplication filter .client_id = <str>$0"
	return database.Client.Execute(database.Context, query, clientId)
}

func GetOAuth2ClientApplication(clientID string) (datatypes.OAuthClient, error) {
	var oauthClient datatypes.OAuthClient
	query := `SELECT OAuthApplication {
	client_id,
	client_secret,
	client_name,
	client_type,
	redirect_uris,
	grant_types,
	scope,
	client_owner: {
		id
	},
	client_description,
	client_homepage_url,
	client_logo_url,
	client_tos_url,
	client_privacy_url,
	client_registration_date,
	client_status } filter .client_id = <str>$0 LIMIT 1`
	return oauthClient, database.Client.QuerySingle(database.Context, query, &oauthClient, clientID)
}

func CreateNewOAuth2AuthorizationCode(authorizationCode datatypes.OAuthAuthorizationCode) error {
	query := `
		INSERT Authcode {
			code := <str>$0,
			application := <OAuthApplication>$1,
			account := <Account>$2,
			requested_scope := <array<str>>$3,
			granted_scope := <array<str>>$4,
			expires_at := <datetime>$5,
			consented := <bool>$6,
			redirect_uri := <str>$7,
		}
	`
	return database.Client.Execute(database.Context, query,
		authorizationCode.Code,
		authorizationCode.Application.ID,
		authorizationCode.Account.Id,
		authorizationCode.RequestedScope,
		authorizationCode.GrantedScope,
		authorizationCode.ExpiresAt,
		authorizationCode.Consented,
		authorizationCode.RedirectURI,
	)
}

func GetOAuth2CleintApplicationAndUserAccount(clientID string, accountID edgedb.UUID) (datatypes.OAuthClient, datatypes.Account, error) {
	var result struct {
		OAuthClient datatypes.OAuthClient `edgedb:"oauth_client"`
		Account     datatypes.Account     `edgedb:"account"`
	}

	query := `SELECT (
    	oauth_client := (SELECT OAuthApplication {
		id,
		client_id,
		client_secret,
		client_name,
		client_type,
		redirect_uris,
		grant_types,
		scope,
		client_description,
		client_homepage_url,
		client_logo_url,
		client_tos_url,
		client_privacy_url,
		client_registration_date,
		client_status } filter .client_id = <str>$0 LIMIT 1),
		account := (SELECT Account filter .id = <uuid>$1 LIMIT 1)
		)`

	return result.OAuthClient, result.Account, database.Client.QuerySingle(
		database.Context,
		query,
		&result,
		clientID,
		accountID,
	)
}

func GetOAuth2AuthorizationCode(code string) (datatypes.OAuthAuthorizationCode, error) {
	var authorizationCode datatypes.OAuthAuthorizationCode
	query := `SELECT Authcode {
	code,
	application: {
		id,
		client_id,
		client_secret,
		client_name,
		client_type,
		redirect_uris,
		grant_types,
		scope,
		client_description,
		client_homepage_url,
		client_logo_url,
		client_tos_url,
		client_privacy_url,
		client_registration_date,
		client_status
	},
	account: {
		id,
		username,
		email,
		otp_secret,
		otp_state
	},
	requested_scope,
	granted_scope,
	expires_at,
	consented,
	redirect_uri
	} filter .code = <str>$0 LIMIT 1`
	return authorizationCode, database.Client.QuerySingle(database.Context, query, &authorizationCode, code)
}

func DeleteOAuth2AuthorizationCode(code string) error {
	query := "DELETE Authcode filter .code = <str>$0"
	return database.Client.Execute(database.Context, query, code)
}

func GetRefreshToken(value string) (datatypes.Token, error) {
	var refreshToken datatypes.Token
	query := `SELECT Token {
		id,
		value,
		scope,
		expires_at,
		revoked,
		variant,
		account: {
			id
		}} filter .value = <str>$0 LIMIT 1`
	return refreshToken, database.Client.QuerySingle(database.Context, query, &refreshToken, value)
}

func DeleteRefreshToken(id edgedb.UUID) error {
	query := "DELETE Token filter .id = <uuid>$0"
	return database.Client.Execute(database.Context, query, id)
}

func DeleteTokens(ids []edgedb.UUID) error {
	query := "DELETE Token filter .id IN array_unpack(<array<str>>$0)"
	return database.Client.Execute(database.Context, query, ids)
}

func DeleteTokensByValue(tokens []string) error {
	query := "DELETE Token filter .value IN array_unpack(<array<str>>$0)"
	return database.Client.Execute(database.Context, query, tokens)
}
