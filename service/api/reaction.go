package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/visione2604/WASA-Exam/service/api/reqcontext"
	"github.com/visione2604/WASA-Exam/service/components/requests"
)

// addReaction adds a reaction to a message
func (rt *_router) addReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get IDs from URL
	conversationID := ps.ByName("conversationId")
	messageID := ps.ByName("messageId")

	if conversationID == "" || messageID == "" {
		http.Error(w, "Missing conversation or message ID", http.StatusBadRequest)
		return
	}

	// Get authenticated user
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request
	var req requests.AddReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode reaction request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if !req.IsValid() {
		http.Error(w, "Invalid reaction type", http.StatusBadRequest)
		return
	}

	// Add reaction
	if err := rt.db.AddReactionToMessage(messageID, userID, req.Type); err != nil {
		ctx.Logger.WithError(err).Error("Failed to add reaction")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// removeReaction removes a reaction from a message
func (rt *_router) removeReaction(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get IDs from URL
	conversationID := ps.ByName("conversationId")
	messageID := ps.ByName("messageId")

	if conversationID == "" || messageID == "" {
		http.Error(w, "Missing conversation or message ID", http.StatusBadRequest)
		return
	}

	// Get authenticated user
	userID, err := rt.getAuthenticatedUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body to get reaction type
	var req struct {
		Type string `json:"type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ctx.Logger.WithError(err).Error("Failed to decode reaction removal request")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate reaction type
	validTypes := map[string]bool{
		"like":       true,
		"heart":      true,
		"laugh":      true,
		"sad_face":   true,
		"angry_face": true,
	}
	if !validTypes[req.Type] {
		http.Error(w, "Invalid reaction type", http.StatusBadRequest)
		return
	}

	// Remove reaction
	if err := rt.db.DeleteReactionFromMessage(messageID, userID, req.Type); err != nil {
		ctx.Logger.WithError(err).Error("Failed to remove reaction")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
