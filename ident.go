package provide

import "fmt"

// Ident client
type Ident struct {
	APIClient
}

// InitIdent convenience method
func InitIdent(token *string) *Ident {
	return &Ident{
		APIClient{
			Host:   "ident.provide.services",
			Path:   "api/v1",
			Scheme: "https",
			Token:  token,
		},
	}
}

// Authenticate a user by email address and password, returning a newly-authorized API token
func Authenticate(email, passwd string) (int, interface{}, error) {
	prvd := InitIdent(nil)
	return prvd.post("authenticate", map[string]interface{}{
		"email":    email,
		"password": passwd,
	})
}

// CreateApplication on behalf of the given API token
func CreateApplication(token string, params map[string]interface{}) (int, interface{}, error) {
	return InitIdent(stringOrNil(token)).post("applications", params)
}

// UpdateApplication using the given API token, application id and params
func UpdateApplication(token, applicationID string, params map[string]interface{}) (int, interface{}, error) {
	uri := fmt.Sprintf("applications/%s", applicationID)
	return InitIdent(stringOrNil(token)).put(uri, params)
}

// ListApplications retrieves a paginated list of applications scoped to the given API token
func ListApplications(token string, params map[string]interface{}) (int, interface{}, error) {
	return InitIdent(stringOrNil(token)).get("applications", params)
}

// GetApplicationDetails retrives application details for the given API token and application id
func GetApplicationDetails(token, applicationID string, params map[string]interface{}) (int, interface{}, error) {
	uri := fmt.Sprintf("applications/%s", applicationID)
	return InitIdent(stringOrNil(token)).get(uri, params)
}

// ListApplicationTokens retrieves a paginated list of application API tokens
func ListApplicationTokens(token, applicationID string, params map[string]interface{}) (int, interface{}, error) {
	uri := fmt.Sprintf("applications/%s/tokens", applicationID)
	return InitIdent(stringOrNil(token)).get(uri, params)
}

// ListTokens retrieves a paginated list of API tokens scoped to the given API token
func ListTokens(token string, params map[string]interface{}) (int, interface{}, error) {
	return InitIdent(stringOrNil(token)).get("tokens", params)
}

// GetTokenDetails retrieves details for the given API token id
func GetTokenDetails(token, tokenID string, params map[string]interface{}) (int, interface{}, error) {
	uri := fmt.Sprintf("tokens/%s", tokenID)
	return InitIdent(stringOrNil(token)).get(uri, params)
}

// DeleteToken removes a previously authorized API token, effectively deauthorizing future calls using the token
func DeleteToken(token, tokenID string) (int, interface{}, error) {
	uri := fmt.Sprintf("tokens/%s", tokenID)
	return InitIdent(stringOrNil(token)).delete(uri)
}

// CreateUser creates a new user for which API tokens and managed signing identities can be authorized
func CreateUser(token string, params map[string]interface{}) (int, interface{}, error) {
	return InitIdent(stringOrNil(token)).post("users", params)
}

// ListUsers retrieves a paginated list of users scoped to the given API token
func ListUsers(token string, params map[string]interface{}) (int, interface{}, error) {
	return InitIdent(stringOrNil(token)).get("users", params)
}

// GetUserDetails retrieves details for the given user id
func GetUserDetails(token, userID string, params map[string]interface{}) (int, interface{}, error) {
	uri := fmt.Sprintf("users/%s", userID)
	return InitIdent(stringOrNil(token)).get(uri, params)
}

// UpdateUser updates an existing user
func UpdateUser(token, userID string, params map[string]interface{}) (int, interface{}, error) {
	uri := fmt.Sprintf("users/%s", userID)
	return InitIdent(stringOrNil(token)).put(uri, params)
}
