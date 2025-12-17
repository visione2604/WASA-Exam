/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Crucio104/wasa/service/components/schema"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	Ping() error

	// user related
	SearchUserByUsername(username string) ([]schema.User, error)
	GetUserByName(username string) (*schema.User, error)
	GetUserById(id string) (*schema.User, error)
	CreateUser(user *schema.User) error
	UpdateUsername(userID, newUsername string) error
	UpdateUserPhoto(userID string, photo []byte) error

	// conversation related
	GetMyConversations(userID string) ([]*schema.Conversation, error)
	GetConversationByID(userID, conversationID string) (*schema.Conversation, error)
	SearchConversationByName(name string) ([]schema.Conversation, error)
	CreateConversation(conversation *schema.Conversation) error
	GetLastMessageByConversationID(conversationID string) (*schema.Message, error)
	EnsureDirectConversation(userID, peerUserID string) (*schema.Conversation, error)
	GetConversationMembers(conversationID string) ([]schema.User, error)
	// Group specific methods (aligned with Conversation internal logic)
	UpdateGroupName(conversationID, newName string) error
	UpdateGroupPhoto(conversationID string, photo []byte) error
	AddUserToGroup(conversationID, userID string) error
	LeaveGroup(conversationID, userID string) error

	// message related
	SendMessage(message *schema.Message) error
	GetMessagesByConversationID(conversationID string) ([]*schema.Message, error)
	GetMessageByID(messageID string) (*schema.Message, error)
	ForwardMessage(message *schema.Message, userID string) error
	DeleteMessage(conversationID, messageID, userID string) error
	MarkMessageStatus(messageID, userID, status string) error

	// reaction related
	AddReactionToMessage(messageID, userID, reactionType string) error
	DeleteReactionFromMessage(messageID, userID, reactionType string) error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		usersTable := `CREATE TABLE users (
			id TEXT NOT NULL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			photo BLOB
		);`

		// Aligned with Swagger: added is_group and group_photo
		conversationsTable := `CREATE TABLE conversations (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			is_group BOOLEAN NOT NULL CHECK (is_group IN (0, 1)),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			group_photo BLOB
		);`

		conversationMembersTable := `CREATE TABLE conversation_members (
			conversationId TEXT NOT NULL,
			userId TEXT NOT NULL,
			PRIMARY KEY (conversationId, userId),
			FOREIGN KEY (conversationId) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);`

		// Split content into type and value to match schema.MessageContent
		messagesTable := `CREATE TABLE messages (
			id TEXT NOT NULL PRIMARY KEY,
			conversationId TEXT NOT NULL,
			senderId TEXT NOT NULL,
			content_type TEXT NOT NULL, 
			content_value TEXT NOT NULL,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			status TEXT NOT NULL,
			forwarded_from TEXT,
			FOREIGN KEY (conversationId) REFERENCES conversations(id) ON DELETE CASCADE,
			FOREIGN KEY (senderId) REFERENCES users(id) ON DELETE CASCADE
		);`

		// Aligned with Swagger: 'type' instead of 'reaction'
		reactionsTable := `CREATE TABLE reactions (
			messageId TEXT NOT NULL,
			userId TEXT NOT NULL,
			type TEXT NOT NULL,
			PRIMARY KEY (messageId, userId, type),
			FOREIGN KEY (messageId) REFERENCES messages(id) ON DELETE CASCADE,
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);`

		messageStatusTable := `CREATE TABLE message_status (
			messageId TEXT NOT NULL,
			userId TEXT NOT NULL,
			deliveredAt DATETIME,
			readAt DATETIME,
			PRIMARY KEY (messageId, userId),
			FOREIGN KEY (messageId) REFERENCES messages(id) ON DELETE CASCADE,
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);`

		creationQueries := []string{
			usersTable,
			conversationsTable,
			conversationMembersTable,
			messagesTable,
			reactionsTable,
			messageStatusTable,
		}
		for _, stmt := range creationQueries {
			_, err = db.Exec(stmt)
			if err != nil {
				return nil, fmt.Errorf("error creating table: %w", err)
			}
		}
	}

	return &appdbimpl{c: db}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}