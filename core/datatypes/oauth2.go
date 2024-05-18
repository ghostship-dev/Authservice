package datatypes

import (
	stdErrors "errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/edgedb/edgedb-go"
	"github.com/ghostship-dev/authservice/core/scopes"
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

type UpdateOAuth2ClientKeyValueRequest struct {
	ClientID string `json:"client_id"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

func (r *UpdateOAuth2ClientKeyValueRequest) Validate() map[string]string {
	errors := make(map[string]string)
	if len(strings.TrimSpace(r.ClientID)) < 1 {
		errors["client_id"] = "client_id is required"
	}
	if len(strings.TrimSpace(r.Key)) < 1 {
		errors["key"] = "key is required"
	}
	if len(strings.TrimSpace(r.Key)) > 1 && r.Key != "client_name" &&
		r.Key != "client_type" &&
		r.Key != "redirect_uris" &&
		r.Key != "grant_types" &&
		r.Key != "scope" &&
		r.Key != "client_description" &&
		r.Key != "client_homepage_url" &&
		r.Key != "client_logo_url" &&
		r.Key != "client_tos_url" &&
		r.Key != "client_privacy_url" {
		errors["key"] = "key (" + r.Key + ") can not be modified"
	}
	if len(strings.TrimSpace(r.Value)) < 1 {
		errors["value"] = "value is required"
	}
	// Validity checks
	if r.Key == "client_name" && len(strings.TrimSpace(r.Value)) < 4 {
		errors["client_name"] = "client_name is must be at least 4 characters long"
	}
	if r.Key == "redirect_uris" && len(r.Value) <= 0 {
		errors["redirect_uris"] = "redirect_uris is required"
	}
	if r.Key == "redirect_uris" {
		regex, err := regexp.Compile(`^(https?)://[^\s/$.?#].\S*$`)
		if err != nil {
			errors["redirect_uris"] = "internal server error"
		} else {
			strArray := strings.Split(strings.TrimSpace(r.Value), ",")
			for i := range strArray {
				if !regex.MatchString(strArray[i]) {
					errors["redirect_uris_"+strconv.Itoa(i)] = "'" + strArray[i] + "' does not match the redirect url requirements."
				}
			}
		}
	}
	if r.Key == "client_logo_url" {
		regex, err := regexp.Compile(`^(https?)://[^\s/$.?#].\S*$`)
		if err != nil {
			errors["client_logo_url"] = "internal server error"
		} else {
			strArray := strings.Split(strings.TrimSpace(r.Value), ",")
			for i := range strArray {
				if !regex.MatchString(strArray[i]) {
					errors["client_logo_url_"+strconv.Itoa(i)] = "'" + strArray[i] + "' does not match the client logo url requirements."
				}
			}
		}
	}
	if r.Key == "client_privacy_url" {
		regex, err := regexp.Compile(`^(https?)://[^\s/$.?#].\S*$`)
		if err != nil {
			errors["client_privacy_url"] = "internal server error"
		} else {
			strArray := strings.Split(strings.TrimSpace(r.Value), ",")
			for i := range strArray {
				if !regex.MatchString(strArray[i]) {
					errors["client_privacy_url_"+strconv.Itoa(i)] = "'" + strArray[i] + "' does not match the client privacy url requirements."
				}
			}
		}
	}
	if r.Key == "client_tos_url" {
		regex, err := regexp.Compile(`^(https?)://[^\s/$.?#].\S*$`)
		if err != nil {
			errors["redirect_uris"] = "internal server error"
		} else {
			strArray := strings.Split(strings.TrimSpace(r.Value), ",")
			for i := range strArray {
				if !regex.MatchString(strArray[i]) {
					errors["client_tos_url_"+strconv.Itoa(i)] = "'" + strArray[i] + "' does not match the client tos url requirements."
				}
			}
		}
	}
	if r.Key == "grant_types" && len(r.Value) < 1 {
		errors["grant_types"] = "grant_types is required"
	}
	if r.Key == "grant_types" {
		for i := range r.Value {
			strArray := strings.Split(strings.TrimSpace(r.Value), ",")
			grantTypeValue := strings.TrimSpace(strArray[i])
			if grantTypeValue != ImplicitGrant && grantTypeValue != PasswordGrant && grantTypeValue != ClientCredentialsGrant && grantTypeValue != AuthorizationCodeGrant {
				errors["grant_types_"+strconv.Itoa(i)] = "invalid grant type: " + grantTypeValue
			}
		}
	}
	if r.Key == "scope" && len(r.Value) < 1 {
		errors["scope"] = "scope is required"
	}
	if r.Key == "scope" {
		strArray := strings.Split(strings.TrimSpace(r.Value), ",")
		if !scopes.AllScopesAllowed(strArray) {
			for i, scope := range strArray {
				if !scopes.IsScopeAllowed(scope) {
					errors["scope_"+strconv.Itoa(i)] = "scope '" + scope + "' is not allowed"
				}
			}
		}
	}
	return errors
}

func (r *UpdateOAuth2ClientKeyValueRequest) GetKeyType() (string, error) {
	if r.Key == "client_name" || r.Key == "client_type" || r.Key == "client_description" || r.Key == "client_homepage_url" || r.Key == "client_logo_url" || r.Key == "client_tos_url" || r.Key == "client_privacy_url" {
		return "<str>", nil
	}
	if r.Key == "redirect_uris" || r.Key == "grant_types" || r.Key == "scope" {
		return "<array<str>>", nil
	}
	return "", stdErrors.New("failed to parse key: " + r.Key)
}

type DeleteOAuth2ClientRequest struct {
	ClientID string `json:"client_id"`
}

func (r *DeleteOAuth2ClientRequest) Validate() map[string]string {
	var errors map[string]string = make(map[string]string)
	if r.ClientID == "" {
		errors["client_id"] = "client_id is required"
	}
	return errors
}

type IntrospectOAuth2TokenRequest struct {
	Token                 string `json:"token"`
	CheckIfTokenIsRevoked bool   `json:"check_if_token_is_revoked"`
}

func (r *IntrospectOAuth2TokenRequest) Validate() map[string]string {
	var errors map[string]string = make(map[string]string)
	if r.Token == "" {
		errors["token"] = "token is required"
	}
	return errors
}
