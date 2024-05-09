package queries

import (
	"context"
	"errors"
	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/database"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func GetPasswordByEmail(email string) (datatypes.Password, error) {
	var password datatypes.Password
	query := "SELECT Password{password, failed_attempts, account: { id }} filter .email = <str>$0 LIMIT 1"
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
