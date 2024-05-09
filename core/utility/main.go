package utility

import (
	"errors"
	"fmt"
	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/golang-jwt/jwt"
	"net/http"
	"os"
	"strings"
	"time"
)

func NewAccessToken(account datatypes.Account, expires time.Time, scope []string) (datatypes.Token, error) {
	tokenString, err := GenerateJWT(expires, account, scope, "access_token")
	if err != nil {
		fmt.Println(err)
		return datatypes.Token{}, err
	}
	return datatypes.Token{
		Variant:   "access_token",
		Value:     tokenString,
		Scope:     scope,
		Account:   account,
		Revoked:   false,
		ExpiresAt: expires,
	}, nil
}

func NewRefreshToken(account datatypes.Account, expires time.Time, scope []string) (datatypes.Token, error) {
	tokenString, err := GenerateJWT(expires, account, scope, "refresh_token")
	if err != nil {
		return datatypes.Token{}, err
	}
	return datatypes.Token{
		Variant:   "refresh_token",
		Value:     tokenString,
		Scope:     scope,
		Account:   account,
		Revoked:   false,
		ExpiresAt: expires,
	}, nil
}

func GenerateJWT(expires time.Time, account datatypes.Account, scope []string, tokenVariant string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expires
	claims["account_id"] = account.Id
	claims["scope"] = scope
	claims["variant"] = tokenVariant

	key := []byte(os.Getenv("JWT_SECRET_KEY"))

	return token.SignedString(key)
}

func GetBearerTokenFromHeader(h *http.Header) (string, error) {
	value := strings.TrimSpace(strings.Replace(h.Get("Authorization"), "Bearer", "", 1))
	if value == "" {
		return "", errors.New("no bearer token found")
	}
	return value, nil
}
