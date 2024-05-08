package utility

import (
	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/datatypes"
	gonanoid "github.com/matoous/go-nanoid"
)

func NewAccessToken(account edgedb.UUID, scope []string) (datatypes.Token, error) {
	tokenValue, err := gonanoid.ID(40)
	if err != nil {
		return datatypes.Token{}, err
	}
	return datatypes.Token{
		Variant: "access_token",
		Value:   tokenValue,
		Scope:   scope,
		Account: account,
		Revoked: false,
	}, nil
}

func NewRefreshToken(account edgedb.UUID, scope []string) (datatypes.Token, error) {
	tokenValue, err := gonanoid.ID(40)
	if err != nil {
		return datatypes.Token{}, err
	}
	return datatypes.Token{
		Variant: "refresh_token",
		Value:   tokenValue,
		Scope:   scope,
		Account: account,
		Revoked: false,
	}, nil
}
