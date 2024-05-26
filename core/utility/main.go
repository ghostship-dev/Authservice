package utility

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ghostship-dev/authservice/core/datatypes"
	"github.com/golang-jwt/jwt/v5"
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
	claims["exp"] = expires.Unix()
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

func GetClientSecretFromHeader(h *http.Header) (string, error) {
	value := strings.TrimSpace(strings.Replace(h.Get("Authorization"), "Basic", "", 1))
	if value == "" {
		return "", errors.New("no client secret found")
	}
	decodedValue, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	credentials := strings.SplitN(string(decodedValue), ":", 2)
	return string(credentials[1]), nil
}

func GetClientIDFromHeader(h *http.Header) (string, error) {
	value := strings.TrimSpace(strings.Replace(h.Get("Authorization"), "Basic", "", 1))
	if value == "" {
		return "", errors.New("no client id found")
	}
	decodedValue, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	credentials := strings.SplitN(string(decodedValue), ":", 2)
	return string(credentials[0]), nil
}
