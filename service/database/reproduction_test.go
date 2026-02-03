package database

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/visione2604/WASA-Exam/service/components/schema"
)

func TestEnsureDirectConversationReusesExisting(t *testing.T) {
	// 1. Setup in-memory DB
	rawDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open db: %v", err)
	}
	defer rawDB.Close()

	// Initialize schema using the existing New function logic
	db, err := New(rawDB)
	if err != nil {
		t.Fatalf("Failed to init database: %v", err)
	}

	// 2. Create two users
	u1 := &schema.User{ID: "u1", Username: "user1"}
	u2 := &schema.User{ID: "u2", Username: "user2"}

	if _, err := db.CreateUser(u1); err != nil {
		t.Fatalf("Failed to create user 1: %v", err)
	}
	if _, err := db.CreateUser(u2); err != nil {
		t.Fatalf("Failed to create user 2: %v", err)
	}

	// 3. Ensure conversation first time
	conv1, err := db.EnsureDirectConversation(u1.ID, u2.ID)
	if err != nil {
		t.Fatalf("EnsureDirectConversation 1 failed: %v", err)
	}
	if conv1 == nil {
		t.Fatal("Expected conversation, got nil")
	}

	// 4. Ensure conversation second time
	conv2, err := db.EnsureDirectConversation(u1.ID, u2.ID)
	if err != nil {
		t.Fatalf("EnsureDirectConversation 2 failed: %v", err)
	}
	if conv2 == nil {
		t.Fatal("Expected conversation, got nil")
	}

	// 5. Assert IDs are the same
	if conv1.ID != conv2.ID {
		t.Errorf("Expected conversation IDs to be equal. First: %s, Second: %s. A new conversation was created instead of reusing the existing one.", conv1.ID, conv2.ID)
	}
}
