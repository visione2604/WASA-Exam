package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/visione2604/WASA-Exam/service/components/schema"
	"github.com/visione2604/WASA-Exam/service/database"
)

func main() {

	dbConn, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	db, err := database.New(dbConn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping failed:", err)
	}
	fmt.Println("Database ready ✔️")

	// --- USERS ---
	aliceID, _ := db.CreateUser(&schema.User{Username: "alice"})
	bobID, _ := db.CreateUser(&schema.User{Username: "bob"})
	charlieID, _ := db.CreateUser(&schema.User{Username: "charlie"})
	fmt.Println("Users created:", aliceID, bobID, charlieID)

	// --- DIRECT CONVERSATIONS ---
	directAB, _ := db.EnsureDirectConversation(aliceID, bobID)
	directAC, _ := db.EnsureDirectConversation(aliceID, charlieID)

	fmt.Println("Direct conv A-B:", directAB.ID)
	fmt.Println("Direct conv A-C:", directAC.ID)

	// --- GROUP CONVERSATION ---
	group := &schema.Conversation{
		Name:    "Project Team",
		IsGroup: true,
		Participants: []schema.User{
			{ID: aliceID},
			{ID: bobID},
			{ID: charlieID},
		},
	}
	_ = db.CreateConversation(group)
	fmt.Println("Group created:", group.ID)

	// --- SEND MESSAGE ---
	msg := &schema.Message{
		Sender: schema.User{ID: aliceID},
		Content: schema.MessageContent{
			Type:  "text",
			Value: "Hello team!",
		},
		Timestamp:     time.Now(),
		MessageStatus: "sent",
	}
	_ = db.SendMessage(msg, group.ID)
	fmt.Println("Message sent:", msg.ID)

	// --- LIST MESSAGES ---
	msgs, _ := db.GetMessagesByConversationID(group.ID)
	fmt.Println("Messages in group:", len(msgs))
	// --- GET Last MESSAGE ---
	lastMsg, _ := db.GetLastMessageByConversationID(group.ID)
	fmt.Println("Last message ID:", lastMsg.ID)

	// --- FORWARD MESSAGE ---
	forward := &schema.Message{
		Sender: schema.User{ID: bobID},
		Content: schema.MessageContent{
			Type:  "text",
			Value: msgs[0].Content.Value,
		},
		ForwardedFrom: msgs[0].ID,
		MessageStatus: "sent",
		Timestamp:     time.Now(),
	}
	_ = db.ForwardMessage(forward, bobID)
	fmt.Println("Message forwarded")

	// --- REACTIONS ---
	_ = db.AddReactionToMessage(msgs[0].ID, charlieID, "👍")
	fmt.Println("Reaction added")

	_ = db.DeleteReactionFromMessage(msgs[0].ID, charlieID, "👍")
	fmt.Println("Reaction removed")

	// --- MEMBERS ---
	members, _ := db.GetConversationMembers(group.ID)
	fmt.Println("Group members:", len(members))

	// --- DELETE MESSAGE ---
	err = db.DeleteMessage(group.ID, msgs[0].ID, aliceID)
	if err != nil {
		fmt.Println("Delete failed:", err)
	} else {
		fmt.Println("Message deleted")
	}

	// --- MY CONVERSATIONS ---
	convs, _ := db.GetMyConversations(aliceID)
	fmt.Println("Alice conversations:", len(convs))

	// --- FETCH FULL CONVERSATION ---
	full, _ := db.GetConversationByID(aliceID, group.ID)
	fmt.Println("Loaded conversation:", full.Name)

	fmt.Println("\n🎉 Test run finished successfully")
}
