package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/visione2604/WASA-Exam/service/components/schema"
)

// SendMessage inserts a new message in a conversation
func (db *appdbimpl) SendMessage(msg *schema.Message, conversationID string) error {
	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	msg.ID = uuid.Must(uuid.NewV4()).String()
	msg.Timestamp = time.Now().UTC()

	_, err := db.c.Exec(`
		INSERT INTO messages (
			id, conversationId, senderId,
			content_type, content_value,
			timestamp, status, forwarded_from
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		msg.ID,
		conversationID,
		msg.Sender.ID,
		msg.Content.Type,
		msg.Content.Value,
		msg.Timestamp,
		msg.MessageStatus,
		msg.ForwardedFrom,
	)

	return err
}

// GetMessagesByConversationID returns all messages for a conversation
// GetMessagesByConversationID returns all messages for a conversation
// GetMessagesByConversationID returns all messages for a conversation
func (db *appdbimpl) GetMessagesByConversationID(conversationID string, currentUserID string) ([]*schema.Message, error) {
	rows, err := db.c.Query(`
		SELECT 
			m.id, 
			m.senderId, 
			m.content_type, 
			m.content_value,
			m.timestamp, 
			m.forwarded_from,
			COALESCE(
				(SELECT 
					CASE 
						WHEN COUNT(CASE WHEN ms.readAt IS NULL THEN 1 END) = 0 THEN 'read'
						WHEN COUNT(CASE WHEN ms.deliveredAt IS NULL THEN 1 END) = 0 THEN 'delivered'
						ELSE 'sent'
					END
				FROM conversation_members cm
				LEFT JOIN message_status ms ON ms.messageId = m.id AND ms.userId = cm.userId
				WHERE cm.conversationId = ? AND cm.userId != m.senderId
				),
				'sent'
			) as message_status
		FROM messages m
		WHERE m.conversationId = ?
		ORDER BY m.timestamp ASC
	`, conversationID, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*schema.Message

	for rows.Next() {
		var m schema.Message
		var senderID string
		var ts string

		if err := rows.Scan(
			&m.ID,
			&senderID,
			&m.Content.Type,
			&m.Content.Value,
			&ts,
			&m.ForwardedFrom,
			&m.MessageStatus,
		); err != nil {
			return nil, err
		}

		// parse timestamp
		if t, perr := time.Parse(time.RFC3339, ts); perr == nil {
			m.Timestamp = t
		}

		// load sender info
		_ = db.c.QueryRow(`
			SELECT id, username, photo
			FROM users WHERE id = ?
		`, senderID).Scan(&m.Sender.ID, &m.Sender.Username, &m.Sender.Photo)

		// attach reactions
		m.Reactions, _ = db.GetReactionsForMessage(m.ID)

		msgs = append(msgs, &m)
	}

	return msgs, rows.Err()
}

// GetMessageByID returns a single message
// GetMessageByID returns a single message
func (db *appdbimpl) GetMessageByID(messageID string) (*schema.Message, error) {
	var m schema.Message
	var senderID string
	var ts string
	var conversationID string

	// Prima ottieni il conversationID
	err := db.c.QueryRow(`
		SELECT conversationId FROM messages WHERE id = ?
	`, messageID).Scan(&conversationID)
	if err != nil {
		return nil, fmt.Errorf("message not found")
	}

	// Poi ottieni il messaggio con lo status calcolato
	err = db.c.QueryRow(`
		SELECT 
			m.id, 
			m.senderId, 
			m.content_type, 
			m.content_value,
			m.timestamp, 
			m.forwarded_from,
			COALESCE(
				(SELECT 
					CASE 
						WHEN COUNT(CASE WHEN ms.readAt IS NULL THEN 1 END) = 0 THEN 'read'
						WHEN COUNT(CASE WHEN ms.deliveredAt IS NULL THEN 1 END) = 0 THEN 'delivered'
						ELSE 'sent'
					END
				FROM conversation_members cm
				LEFT JOIN message_status ms ON ms.messageId = m.id AND ms.userId = cm.userId
				WHERE cm.conversationId = ? AND cm.userId != m.senderId
				),
				'sent'
			) as message_status
		FROM messages m
		WHERE m.id = ?
	`, conversationID, messageID).Scan(
		&m.ID,
		&senderID,
		&m.Content.Type,
		&m.Content.Value,
		&ts,
		&m.ForwardedFrom,
		&m.MessageStatus,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("message not found")
		}
		return nil, err
	}

	if t, perr := time.Parse(time.RFC3339, ts); perr == nil {
		m.Timestamp = t
	}

	// load sender
	_ = db.c.QueryRow(`
		SELECT id, username, photo
		FROM users WHERE id = ?
	`, senderID).Scan(&m.Sender.ID, &m.Sender.Username, &m.Sender.Photo)

	// attach reactions
	m.Reactions, _ = db.GetReactionsForMessage(m.ID)

	return &m, nil
}

// DeleteMessage deletes a message if the requesting user is the sender
func (db *appdbimpl) DeleteMessage(conversationID, messageID, userID string) error {
	// Check the message exists and belongs to the user in that conversation
	var senderID string
	err := db.c.QueryRow(`
		SELECT senderId
		FROM messages
		WHERE id = ? AND conversationId = ?
	`, messageID, conversationID).Scan(&senderID)

	if err == sql.ErrNoRows {
		return fmt.Errorf("message not found")
	}
	if err != nil {
		return fmt.Errorf("failed to verify message owner: %w", err)
	}

	// Only the sender can delete their message
	if senderID != userID {
		return fmt.Errorf("unauthorized to delete message")
	}

	// Perform deletion
	_, err = db.c.Exec(`
		DELETE FROM messages
		WHERE id = ?
	`, messageID)

	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

// ForwardMessage clones an existing message into another conversation
func (db *appdbimpl) ForwardMessage(msg *schema.Message, userID string) error {
	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	msg.ID = uuid.Must(uuid.NewV4()).String()
	msg.Timestamp = time.Now().UTC()
	msg.ForwardedFrom = userID

	_, err := db.c.Exec(`
		INSERT INTO messages (
			id, conversationId, senderId,
			content_type, content_value,
			timestamp, status, forwarded_from
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		msg.ID,
		msg.ConversationID,
		msg.Sender.ID,
		msg.Content.Type,
		msg.Content.Value,
		msg.Timestamp,
		msg.MessageStatus,
		msg.ForwardedFrom,
	)

	return err
}

// MarkMessageStatus updates delivery / read receipts
func (db *appdbimpl) MarkMessageStatus(messageID, userID, status string) error {
	switch status {
	case "delivered":
		_, err := db.c.Exec(`
			INSERT INTO message_status (messageId, userId, deliveredAt)
			VALUES (?, ?, datetime('now'))
			ON CONFLICT(messageId, userId)
			DO UPDATE SET deliveredAt = datetime('now')
		`, messageID, userID)
		return err

	case "read":
		_, err := db.c.Exec(`
			UPDATE message_status
			SET readAt = datetime('now')
			WHERE messageId = ? AND userId = ?
		`, messageID, userID)
		return err
	}

	return fmt.Errorf("invalid message status: %s", status)
}
