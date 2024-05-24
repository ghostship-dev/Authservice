package responses

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/ghostship-dev/authservice/core/datatypes"
)

type genericOAuth2ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func makeResponse(statusCode int, message string) error {
	response := genericOAuth2ErrorResponse{
		Error:   true,
		Message: message,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return InternalServerErrorResponse()
	}
	return datatypes.NewRequestError(statusCode, string(jsonResponse))
}

func OAuth2ApplicationNotFoundResponse() error {
	return makeResponse(http.StatusBadRequest, "oauth2 application not found")
}

func OAuth2RedirectURIDoesNotMatch() error {
	return makeResponse(http.StatusBadRequest, "redirect_uri does not match one of the registered redirect URIs for this client")
}

func OAuth2ScopeIsRequired() error {
	return makeResponse(http.StatusBadRequest, "scope is required")
}

func OAuth2InvalidScope(scope []string) error {
	return makeResponse(http.StatusBadRequest, "invalid scope '"+strings.Join(scope, ", ")+"'")
}

func OAuth2UserNotFoundResponse() error {
	return makeResponse(http.StatusBadRequest, "user not found")
}

func ReturnRedirectResponseToConsentPage(w http.ResponseWriter, r *http.Request, authCode datatypes.OAuthAuthorizationCode) error {
	http.Redirect(w, r, os.Getenv("OAuth2_ConsentPage_URI")+"?code="+authCode.Code, http.StatusSeeOther)
	return nil
}
