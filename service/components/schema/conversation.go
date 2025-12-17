package schema

type Conversation struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	GroupPhoto   []byte    `json:"groupPhoto,omitempty"` 
	IsGroup      bool      `json:"isGroup"` 
	
	Participants []User    `json:"participants"` 
	
	LastMessage  *Message  `json:"lastMessage,omitempty"` 
	
	Messages     []Message `json:"messages,omitempty"`
}

type SetGroupNameRequest struct {
	Name string `json:"name"`
}

type SetGroupPhotoRequest struct {
	Photo []byte `json:"photo"`
}