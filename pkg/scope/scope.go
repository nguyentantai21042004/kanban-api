package scope

import (
	"encoding/base64"
	"encoding/json"
)

type Scope struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// NewScope creates a new scope.
func NewScope(payload Payload) Scope {
	return Scope{
		UserID: payload.UserID,
		Email:  payload.Email,
	}
}

func CreateScopeHeader(scope Scope) (string, error) {
	// Marshal the scope data to JSON
	jsonData, err := json.Marshal(scope)
	if err != nil {
		return "", err
	}

	// Encode the JSON data as Base64
	base64Data := base64.StdEncoding.EncodeToString(jsonData)
	return base64Data, nil
}

func ParseScopeHeader(scopeHeader string) (Scope, error) {
	// Decode the Base64 data
	jsonData, err := base64.StdEncoding.DecodeString(scopeHeader)
	if err != nil {
		return Scope{}, err
	}

	// Unmarshal the JSON data
	var scope Scope
	err = json.Unmarshal(jsonData, &scope)
	if err != nil {
		return Scope{}, err
	}

	return scope, nil
}

func (m implManager) VerifyScope(scopeHeader string) (Scope, error) {
	scope, err := ParseScopeHeader(scopeHeader)
	if err != nil {
		return Scope{}, err
	}

	return scope, nil
}
