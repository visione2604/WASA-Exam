package schema

// this is a struct that represents a user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Photo    []byte `json:"photo,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	User  User   // it will return the user info
	Token string `json:"token"`
}

// UsernameUpdateRequest is used when we do a username change
type UsernameUpdateRequest struct {
	Username string `json:"username"`
}

// UsernameUpdateResponse returns the updated user info
type UsernameUpdateResponse = User // Alias, same structure as User

// ProfilePhotoUpdateRequest is used when we do an update profile photo
type ProfilePhotoUpdateRequest struct {
	Photo []byte `json:"photo"`
}

// ProfilePhotoUpdateResponse returns the updated photo
type ProfilePhotoUpdateResponse struct {
	Photo []byte `json:"photo"`
}
