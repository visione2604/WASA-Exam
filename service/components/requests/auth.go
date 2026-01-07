package requests

import (
	"regexp"
)

// LoginRequest contains the username submitted during login
type LoginRequest struct {
	Name string `json:"name"`
}

// IsValid checks if the username meets required format constraints
// According to swagger: minLength: 3, maxLength: 20
func (r *LoginRequest) IsValid() bool {
	if len(r.Name) < 3 || len(r.Name) > 20 {
		return false
	}
	// Accepts alphanumeric usernames with underscores and spaces
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_ ]+$`, r.Name)
	return match
}
