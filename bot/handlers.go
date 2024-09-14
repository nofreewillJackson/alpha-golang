// bot/handlers.go
package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/yourusername/myproject/common"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Store the message in the database
	msg := common.Message{
		Content:    m.Content,
		AuthorID:   m.Author.ID,
		ChannelID:  m.ChannelID,
		Timestamp:  time.Now(),
		Summarized: false,
	}

	_, err := dbpool.Exec(context.Background(),
		`INSERT INTO messages (content, author_id, channel_id, timestamp, summarized)
         VALUES ($1, $2, $3, $4, $5)`,
		msg.Content, msg.AuthorID, msg.ChannelID, msg.Timestamp, msg.Summarized)

	if err != nil {
		log.Printf("Error inserting message: %v\n", err)
	}

	// Check if threshold is reached and summarize
	checkAndSummarizeMessages()
}

func checkAndSummarizeMessages() {
	// Define threshold (e.g., 5000 characters)
	const threshold = 5000

	// Fetch unsummarized messages
	rows, err := dbpool.Query(context.Background(),
		`SELECT id, content FROM messages WHERE summarized = false`)
	if err != nil {
		log.Printf("Error fetching messages: %v\n", err)
		return
	}
	defer rows.Close()

	var contents []string
	var messageIDs []int

	var totalLength int
	for rows.Next() {
		var id int
		var content string
		err := rows.Scan(&id, &content)
		if err != nil {
			log.Printf("Error scanning message: %v\n", err)
			continue
		}

		totalLength += len(content)
		contents = append(contents, content)
		messageIDs = append(messageIDs, id)

		if totalLength >= threshold {
			break
		}
	}

	if len(contents) == 0 {
		return
	}

	// Summarize using OpenAI
	combinedContent := strings.Join(contents, "\n")
	summary, err := summarizeContent(combinedContent)
	if err != nil {
		log.Printf("Error summarizing content: %v\n", err)
		return
	}

	// Store the summary
	_, err = dbpool.Exec(context.Background(),
		`INSERT INTO summaries (summary) VALUES ($1)`, summary)
	if err != nil {
		log.Printf("Error inserting summary: %v\n", err)
		return
	}

	// Mark messages as summarized
	_, err = dbpool.Exec(context.Background(),
		`UPDATE messages SET summarized = true WHERE id = ANY($1)`, messageIDs)
	if err != nil {
		log.Printf("Error updating messages: %v\n", err)
	}
}
