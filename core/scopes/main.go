package scopes

var AllowedScopes = map[string]bool{
	"profile":       true,
	"email":         true,
	"account_write": true,
	"account_read":  true,
}

func IsScopeAllowed(scope string) bool {
	return AllowedScopes[scope]
}

func AllScopesAllowed(scope []string) bool {
	for _, scope := range scope {
		if !IsScopeAllowed(scope) {
			return false
		}
	}
	return true
}
