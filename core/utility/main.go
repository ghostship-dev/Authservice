package utility

import (
	"fmt"
	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func NewAccessToken(accountId edgedb.UUID, expires time.Time, scope []string) (datatypes.Token, error) {
	tokenString, err := GenerateJWT(expires, accountId, scope, "access_token")
	if err != nil {
		fmt.Println(err)
		return datatypes.Token{}, err
	}
	return datatypes.Token{
		Variant:   "access_token",
		Value:     tokenString,
		Scope:     scope,
		Account:   accountId,
		Revoked:   false,
		ExpiresAt: expires,
	}, nil
}

func NewRefreshToken(accountId edgedb.UUID, expires time.Time, scope []string) (datatypes.Token, error) {
	tokenString, err := GenerateJWT(expires, accountId, scope, "refresh_token")
	if err != nil {
		return datatypes.Token{}, err
	}
	return datatypes.Token{
		Variant:   "refresh_token",
		Value:     tokenString,
		Scope:     scope,
		Account:   accountId,
		Revoked:   false,
		ExpiresAt: expires,
	}, nil
}

func GenerateJWT(expires time.Time, accountId edgedb.UUID, scope []string, tokenVariant string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expires
	claims["account_id"] = accountId
	claims["scope"] = scope
	claims["variant"] = tokenVariant

	key := []byte(os.Getenv("JWT_SECRET_KEY"))

	return token.SignedString(key)
}
