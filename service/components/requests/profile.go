package requests

import (
	"regexp"
)

// UsernameUpdateRequest is used when changing username
type UsernameUpdateRequest struct {
	Name string `json:"name"`
}

// IsValid checks if the new username meets required format constraints
func (r *UsernameUpdateRequest) IsValid() bool {
	if len(r.Name) < 3 || len(r.Name) > 20 {
		return false
	}
	// Accepts alphanumeric usernames with underscores and spaces
	match, _ := regexp.MatchString(`^[a-zA-Z0-9_ ]+$`, r.Name)
	return match
}

// ProfilePhotoUpdateRequest is used when updating profile photo
type ProfilePhotoUpdateRequest struct {
	Photo []byte `json:"photo"`
}

// IsValid checks if the photo data is present and within size limits
func (r *ProfilePhotoUpdateRequest) IsValid() bool {
	// Photo must be present and not exceed 5MB
	return len(r.Photo) > 0 && len(r.Photo) <= (5*1024*1024)
}
