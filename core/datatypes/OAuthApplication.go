package datatypes

import (
	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/scopes"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	AuthorizationCodeGrant string = "authorization_code"
	ImplicitGrant          string = "implicit"
	PasswordGrant          string = "password"
	ClientCredentialsGrant string = "client_credentials"
)

const (
	OAuthApplicationStatusActive    string = "active"
	OAuthApplicationStatusDisabled  string = "disabled"
	OAuthApplicationStatusSuspended string = "suspended"
)

type OAuthClient struct {
	ClientID               string             `json:"client_id" edgedb:"client_id"`
	ClientSecret           string             `json:"client_secret" edgedb:"client_secret"`
	ClientName             string             `json:"client_name" edgedb:"client_name"`
	ClientType             string             `json:"client_type" edgedb:"client_type"`
	RedirectURIs           []string           `json:"redirect_uris" edgedb:"redirect_uris"`
	GrantTypes             []string           `json:"grant_types" edgedb:"grant_types"`
	Scope                  []string           `json:"scope" edgedb:"scope"`
	ClientOwner            Account            `json:"client_owner" edgedb:"client_owner"`
	ClientDescription      edgedb.OptionalStr `json:"client_description" edgedb:"client_description"`
	ClientHomepageUrl      edgedb.OptionalStr `json:"client_homepage_url" edgedb:"client_homepage_url"`
	ClientLogoUrl          edgedb.OptionalStr `json:"client_logo_url" edgedb:"client_logo_url"`
	ClientTosUrl           edgedb.OptionalStr `json:"client_tos_url" edgedb:"client_tos_url"`
	ClientPrivacyUrl       edgedb.OptionalStr `json:"client_privacy_url" edgedb:"client_privacy_url"`
	ClientRegistrationDate time.Time          `json:"client_registration_date" edgedb:"client_registration_date"`
	ClientStatus           string             `json:"client_status" edgedb:"client_status"`
	ClientRateLimits       []byte             `json:"client_rate_limits" edgedb:"client_rate_limits"`
}

type NewOAuthClientRequest struct {
	ClientName        string      `json:"client_name"`
	ClientType        string      `json:"client_type"`
	RedirectUris      []string    `json:"redirect_uris"`
	GrantTypes        []string    `json:"grant_types"`
	Scope             []string    `json:"scope"`
	ClientOwner       edgedb.UUID `json:"client_owner"`
	ClientDescription string      `json:"client_description"`
	ClientHomepageUrl string      `json:"client_homepage_url"`
	ClientLogoUrl     string      `json:"client_logo_url"`
	ClientTosUrl      string      `json:"client_tos_url"`
	ClientPrivacyUrl  string      `json:"client_privacy_url"`
}

func (r *NewOAuthClientRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if len(strings.TrimSpace(r.ClientName)) < 4 {
		errors["client_name"] = "client_name is must be at least 4 characters long"
	}
	if len(strings.TrimSpace(r.ClientType)) < 1 {
		errors["client_type"] = "client_type is required"
	}
	if len(r.RedirectUris) <= 0 {
		errors["redirect_uris"] = "redirect_uris is required"
	} else {
		regex, err := regexp.Compile(`^(https?)://[^\s/$.?#].\S*$`)
		if err != nil {
			errors["redirect_uris"] = "internal server error"
		} else {
			for i := range r.RedirectUris {
				if !regex.MatchString(r.RedirectUris[i]) {
					errors["redirect_uris_"+strconv.Itoa(i)] = "'" + r.RedirectUris[i] + "' does not match the redirect url requirements."
				}
			}
		}
	}
	if len(r.GrantTypes) < 1 {
		errors["grant_types"] = "grant_types is required"
	} else {
		for i := range r.GrantTypes {
			grantTypeValue := strings.TrimSpace(r.GrantTypes[i])
			if grantTypeValue != ImplicitGrant && grantTypeValue != PasswordGrant && grantTypeValue != ClientCredentialsGrant && grantTypeValue != AuthorizationCodeGrant {
				errors["grant_types_"+strconv.Itoa(i)] = "invalid grant type: " + grantTypeValue
			}
		}
	}
	if len(r.Scope) < 1 {
		errors["scope"] = "scope is required"
	} else {
		if !scopes.AllScopesAllowed(r.Scope) {
			for i, scope := range r.Scope {
				if !scopes.IsScopeAllowed(scope) {
					errors["scope_"+strconv.Itoa(i)] = "scope '" + scope + "' is not allowed"
				}
			}
		}
	}
	if len(r.ClientOwner.String()) < 1 {
		errors["client_owner"] = "client_owner is required"
	}
	return errors
}
