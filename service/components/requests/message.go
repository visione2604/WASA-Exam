package requests

// MessageContent represents the content of a message
type MessageContent struct {
	Type  string `json:"type"`  // "text" or "photo"
	Value string `json:"value"` // text content or base64 encoded photo
}

// IsValid checks if the message content is valid
func (c *MessageContent) IsValid() bool {
	if c.Type != "text" && c.Type != "photo" {
		return false
	}

	if c.Type == "text" {
		// Text must be present and not exceed 1000 characters
		return len(c.Value) >= 1 && len(c.Value) <= 1000
	}

	if c.Type == "photo" {
		// Photo must be present and not exceed 5MB (base64 encoded)
		return len(c.Value) >= 1 && len(c.Value) <= (7*1024*1024) // ~5MB after base64 encoding
	}

	return false
}

// SendMessageRequest is used when sending a message to a conversation
type SendMessageRequest struct {
	Content MessageContent `json:"content"`
}

// IsValid checks if the message request is valid
func (r *SendMessageRequest) IsValid() bool {
	return r.Content.IsValid()
}

// DeleteMessageRequest is used when deleting a message
type DeleteMessageRequest struct {
	ConversationID string `json:"conversationId"`
	MessageID      string `json:"messageId"`
}

// IsValid checks if the delete request has valid IDs
func (r *DeleteMessageRequest) IsValid() bool {
	return len(r.ConversationID) > 0 && len(r.ConversationID) <= 36 &&
		len(r.MessageID) > 0 && len(r.MessageID) <= 36
}

// ForwardMessageRequest is used when forwarding a message to another conversation
type ForwardMessageRequest struct {
	TargetConversationID string `json:"targetConversationId"`
}

// IsValid checks if the forward request has a valid target conversation ID
func (r *ForwardMessageRequest) IsValid() bool {
	return len(r.TargetConversationID) > 0 && len(r.TargetConversationID) <= 36
}
