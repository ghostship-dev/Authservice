package queries

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/responses"
	"golang.org/x/crypto/bcrypt"
)

func connectToEdgeDB() (*edgedb.Client, error) {
	os.Setenv("EDGEDB_INSTANCE", "Ghostship")
	ctx := context.Background()
	client, err := edgedb.CreateClient(ctx, edgedb.Options{})
	return client, err
}

type EdgeDBQueries struct {
	client  *edgedb.Client
	context context.Context
}

func NewEdgeDBQueryImplementation() *EdgeDBQueries {
	client, err := connectToEdgeDB()
	if err != nil {
		panic(err)
	}
	return &EdgeDBQueries{
		client:  client,
		context: context.Background(),
	}
}

func (edb *EdgeDBQueries) GetPasswordByEmail(email string) (datatypes.Password, error) {
	var password datatypes.Password
	query := "SELECT Password{password, failed_attempts, account: { id, otp_secret, otp_state }} filter .email = <str>$0 LIMIT 1"
	err := edb.client.QuerySingle(edb.context, query, &password, email)
	return password, err
}

func (edb *EdgeDBQueries) CreateAccount(email, username, password string) (datatypes.Account, error) {
	var account datatypes.Account
	err := edb.client.Tx(edb.context, func(ctx context.Context, tx *edgedb.Tx) error {
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

func (edb *EdgeDBQueries) GetAccountById(id string) (datatypes.Account, error) {
	var account datatypes.Account
	query := "SELECT Account filter .id = <uuid>$0 LIMIT 1"
	accountId, err := edgedb.ParseUUID(id)
	if err != nil {
		return account, err
	}
	err = edb.client.QuerySingle(edb.context, query, &account, accountId)
	return account, err
}

func (edb *EdgeDBQueries) IncrementFailedPasswordLoginAttempts(email string) error {
	query := "UPDATE Password filter .email = <str>$0 set { failed_attempts := .failed_attempts +1, last_failed_attempt := <datetime>$1 }"
	return edb.client.Execute(edb.context, query, email, time.Now())
}

func (edb *EdgeDBQueries) ResetFailedPasswordLoginAttempts(email string) error {
	query := "UPDATE Password filter .email = <str>$0 set { failed_attempts := 0, last_failed_attempt := {} }"
	return edb.client.Execute(edb.context, query, email)
}

func (edb *EdgeDBQueries) AddNewToken(accountId edgedb.UUID, value, variant string, expiresAt time.Time, scope []string) error {
	query := "INSERT Token { account := <Account>$0, variant := <str>$1, scope := <array<str>>$2, value := <str>$3, revoked := <bool>$4, expires_at := <datetime>$5 }"
	return edb.client.Execute(edb.context, query, accountId, variant, scope, value, false, expiresAt)
}

func (edb *EdgeDBQueries) AddNewTokenPair(accountId edgedb.UUID, accessTokenValue, refreshTokenValue string, accessTokenExpiresAt, refreshTokenExpiresAt time.Time, scope []string) error {
	query := "INSERT Token { account := <Account>$0, variant := <str>$1, scope := <array<str>>$2, value := <str>$3, revoked := <bool>$4, expires_at := <datetime>$5 }; INSERT Token { account := <Account>$6, variant := <str>$7, scope := <array<str>>$8, value := <str>$9, revoked := <bool>$10, expires_at := <datetime>$11 }"
	return edb.client.Execute(edb.context, query, accountId, "access_token", scope, accessTokenValue, false, accessTokenExpiresAt, accountId, "refresh_token", scope, refreshTokenValue, false, refreshTokenExpiresAt)
}

func (edb *EdgeDBQueries) GetToken(tokenValue string) (datatypes.Token, error) {
	var token datatypes.Token
	query := "SELECT Token { value, scope, revoked, variant, expires_at, account: { id, username, otp_secret, otp_state } } filter .value = <str>$0 LIMIT 1"
	return token, edb.client.QuerySingle(edb.context, query, &token, tokenValue)
}

func (edb *EdgeDBQueries) ResetOTP(accountId edgedb.UUID) error {
	query := "UPDATE Account filter .id = <uuid>$0 set { otp_secret := <str>{}, otp_state := <str>'disabled' }"
	return edb.client.Execute(edb.context, query, accountId)
}

func (edb *EdgeDBQueries) SetOTPSecret(accountId edgedb.UUID, otpSecret string) error {
	query := "UPDATE Account filter .id = <uuid>$0 set { otp_secret := <str>$1, otp_state := <str>$2 }"
	return edb.client.Execute(edb.context, query, accountId, otpSecret, "verifying")
}

func (edb *EdgeDBQueries) SetOTPState(accountId edgedb.UUID, otpState string) error {
	if otpState != "disabled" && otpState != "enabled" && otpState != "verifying" {
		return errors.New("allowed otp state: disabled, enabled or verifying")
	}
	query := "UPDATE Account filter .id = <uuid>$0 set { otp_state := <str>$1 }"
	return edb.client.Execute(edb.context, query, accountId, otpState)
}

func (edb *EdgeDBQueries) CreateNewOAuthClientApplication(oauthClient datatypes.OAuthClient) error {
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

	return edb.client.Execute(edb.context, query,
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

func (edb *EdgeDBQueries) UpdateOAuth2ClientApplicationKeyValue(updateRequestData datatypes.UpdateOAuth2ClientKeyValueRequest) error {
	keyType, err := updateRequestData.GetKeyType()
	if err != nil {
		return responses.BadRequestResponse()
	}
	query := "UPDATE OAuthApplication filter .client_id = <str>$0 set { " + updateRequestData.Key + " := " + keyType + "'" + updateRequestData.Value + "' }"
	return edb.client.Execute(edb.context, query, updateRequestData.ClientID)
}

func (edb *EdgeDBQueries) DeleteOAuth2ClientApplication(clientId string) error {
	query := "DELETE OAuthApplication filter .client_id = <str>$0"
	return edb.client.Execute(edb.context, query, clientId)
}

func (edb *EdgeDBQueries) GetOAuth2ClientApplication(clientID string) (datatypes.OAuthClient, error) {
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
	return oauthClient, edb.client.QuerySingle(edb.context, query, &oauthClient, clientID)
}

func (edb *EdgeDBQueries) CreateNewOAuth2AuthorizationCode(authorizationCode datatypes.OAuthAuthorizationCode) error {
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
	return edb.client.Execute(edb.context, query,
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

func (edb *EdgeDBQueries) GetOAuth2ClientApplicationAndUserAccount(clientID string, accountID edgedb.UUID) (datatypes.OAuthClient, datatypes.Account, error) {
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

	return result.OAuthClient, result.Account, edb.client.QuerySingle(
		edb.context,
		query,
		&result,
		clientID,
		accountID,
	)
}

func (edb *EdgeDBQueries) GetOAuth2AuthorizationCode(code string) (datatypes.OAuthAuthorizationCode, error) {
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
	return authorizationCode, edb.client.QuerySingle(edb.context, query, &authorizationCode, code)
}

func (edb *EdgeDBQueries) DeleteOAuth2AuthorizationCode(code string) error {
	query := "DELETE Authcode filter .code = <str>$0"
	return edb.client.Execute(edb.context, query, code)
}

func (edb *EdgeDBQueries) GetRefreshToken(value string) (datatypes.Token, error) {
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
	return refreshToken, edb.client.QuerySingle(edb.context, query, &refreshToken, value)
}

func (edb *EdgeDBQueries) DeleteRefreshToken(id edgedb.UUID) error {
	query := "DELETE Token filter .id = <uuid>$0"
	return edb.client.Execute(edb.context, query, id)
}

func (edb *EdgeDBQueries) DeleteTokens(ids []edgedb.UUID) error {
	query := "DELETE Token filter .id IN array_unpack(<array<str>>$0)"
	return edb.client.Execute(edb.context, query, ids)
}

func (edb *EdgeDBQueries) DeleteTokensByValue(tokens []string) error {
	query := "DELETE Token filter .value IN array_unpack(<array<str>>$0)"
	return edb.client.Execute(edb.context, query, tokens)
}
