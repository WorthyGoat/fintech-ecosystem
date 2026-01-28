package scopes

import (
	"strings"
)

// Predefined API scopes
const (
	// Payment scopes
	PaymentsRead  = "payments:read"
	PaymentsWrite = "payments:write"

	// Ledger scopes
	LedgerRead  = "ledger:read"
	LedgerWrite = "ledger:write"

	// Wildcard scope (full access)
	All = "*"
)

// ValidScopes is the set of all valid scope strings.
var ValidScopes = map[string]bool{
	PaymentsRead:  true,
	PaymentsWrite: true,
	LedgerRead:    true,
	LedgerWrite:   true,
	All:           true,
}

// EndpointScope maps path prefixes and HTTP methods to required scopes.
type EndpointScope struct {
	PathPrefix string
	Method     string // GET, POST, PUT, DELETE, or * for any
	Scope      string
}

// EndpointScopes defines which scope is required for each endpoint.
var EndpointScopes = []EndpointScope{
	// Payments endpoints
	{PathPrefix: "/payments", Method: "GET", Scope: PaymentsRead},
	{PathPrefix: "/payments", Method: "POST", Scope: PaymentsWrite},
	{PathPrefix: "/payments", Method: "PUT", Scope: PaymentsWrite},
	{PathPrefix: "/payments", Method: "DELETE", Scope: PaymentsWrite},

	// Ledger endpoints
	{PathPrefix: "/ledger", Method: "GET", Scope: LedgerRead},
	{PathPrefix: "/ledger", Method: "POST", Scope: LedgerWrite},
	{PathPrefix: "/ledger", Method: "PUT", Scope: LedgerWrite},
	{PathPrefix: "/ledger", Method: "DELETE", Scope: LedgerWrite},
}

// GetRequiredScope returns the scope required for a given path and method.
// Returns empty string if no scope is required (public endpoint).
func GetRequiredScope(path, method string) string {
	for _, es := range EndpointScopes {
		if strings.HasPrefix(path, es.PathPrefix) {
			if es.Method == "*" || es.Method == method {
				return es.Scope
			}
		}
	}
	return ""
}

// HasScope checks if the provided scopes include the required scope.
// Supports wildcard (*) which grants all permissions.
func HasScope(scopes string, required string) bool {
	if required == "" {
		return true // No scope required
	}

	// Parse comma or space separated scopes
	scopeList := ParseScopes(scopes)

	for _, s := range scopeList {
		if s == All || s == required {
			return true
		}
		// Check for prefix match (e.g., "payments:*" matches "payments:read")
		if strings.HasSuffix(s, ":*") {
			prefix := strings.TrimSuffix(s, ":*")
			if strings.HasPrefix(required, prefix+":") {
				return true
			}
		}
	}
	return false
}

// ParseScopes parses a scope string (comma or space separated) into a slice.
func ParseScopes(scopes string) []string {
	if scopes == "" {
		return nil
	}

	// Replace commas with spaces and split
	normalized := strings.ReplaceAll(scopes, ",", " ")
	parts := strings.Fields(normalized)

	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// ValidateScopes checks if all scopes in the string are valid.
func ValidateScopes(scopes string) ([]string, []string) {
	scopeList := ParseScopes(scopes)
	var valid, invalid []string

	for _, s := range scopeList {
		if ValidScopes[s] {
			valid = append(valid, s)
		} else {
			invalid = append(invalid, s)
		}
	}
	return valid, invalid
}

// JoinScopes combines multiple scopes into a space-separated string.
func JoinScopes(scopes []string) string {
	return strings.Join(scopes, " ")
}
