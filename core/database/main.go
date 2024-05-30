package database

import (
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/config"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/ghostship-dev/authservice/core/queries"
)

type DatabaseQueries interface {
	GetPasswordByEmail(email string) (datatypes.Password, error)
	CreateAccount(email, username, password string) (datatypes.Account, error)
	GetAccountById(id string) (datatypes.Account, error)
	IncrementFailedPasswordLoginAttempts(email string) error
	ResetFailedPasswordLoginAttempts(email string) error
	AddNewToken(accountId edgedb.UUID, value, variant string, expiresAt time.Time, scope []string) error
	AddNewTokenPair(accountId edgedb.UUID, accessTokenValue, refreshTokenValue string, accessTokenExpiresAt, refreshTokenExpiresAt time.Time, scope []string) error
	GetToken(tokenValue string) (datatypes.Token, error)
	ResetOTP(accountId edgedb.UUID) error
	SetOTPSecret(accountId edgedb.UUID, otpSecret string) error
	SetOTPState(accountId edgedb.UUID, otpState string) error
	CreateNewOAuthClientApplication(oauthClient datatypes.OAuthClient) error
	UpdateOAuth2ClientApplicationKeyValue(updateRequestData datatypes.UpdateOAuth2ClientKeyValueRequest) error
	DeleteOAuth2ClientApplication(clientId string) error
	GetOAuth2ClientApplication(clientID string) (datatypes.OAuthClient, error)
	CreateNewOAuth2AuthorizationCode(authorizationCode datatypes.OAuthAuthorizationCode) error
	GetOAuth2ClientApplicationAndUserAccount(clientID string, accountID edgedb.UUID) (datatypes.OAuthClient, datatypes.Account, error)
	GetOAuth2AuthorizationCode(code string) (datatypes.OAuthAuthorizationCode, error)
	DeleteOAuth2AuthorizationCode(code string) error
	GetRefreshToken(value string) (datatypes.Token, error)
	DeleteRefreshToken(id edgedb.UUID) error
	DeleteTokens(ids []edgedb.UUID) error
	DeleteTokensByValue(tokens []string) error
}

type Database struct {
	Queries DatabaseQueries
}

func ConnectToSelectedDBDriver(c *config.Config) *Database {
	switch "edgedb" {
	case "edgedb":
		return &Database{Queries: queries.NewEdgeDBQueryImplementation(c)}
	}
	panic("No or invalid database driver selected. Please set the DATABASE_ENGINE environment variable to a valid database driver. (available drivers: edgedb)")
}

var Connection *Database
