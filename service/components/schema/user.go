package schema

type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Photo    []byte `json:"profilePhoto,omitempty"`
}

type LoginRequest struct {
    Username string `json:"name"`
}

type LoginResponse struct {
    User  User   `json:"user"`
    Token string `json:"identifier"`
}

type UsernameUpdateResponse = User

type ProfilePhotoUpdateResponse struct {
    Photo []byte `json:"photo"`
}