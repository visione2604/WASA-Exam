package requests

// GetConversationByIDRequest is used when fetching a single conversation
type GetConversationByIDRequest struct {
	ConversationID string `json:"conversationId"`
}

// IsValid checks if conversation ID is present
func (r *GetConversationByIDRequest) IsValid() bool {
	return len(r.ConversationID) > 0 && len(r.ConversationID) <= 36
}

// GetMyConversationsRequest is used when fetching all user conversations
// No additional fields needed - user ID comes from auth token
type GetMyConversationsRequest struct {
	// Placeholder for future pagination or filters
}

// IsValid always returns true as this request has no body parameters
func (r *GetMyConversationsRequest) IsValid() bool {
	return true
}
