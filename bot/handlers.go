// bot/handlers.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nofreewilljackson/alpha-golang/common"
)

const summarizeChannelID = "1279549793286225930" // Channel ID where summaries will be sent

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Handle command: /summarize
	if strings.HasPrefix(m.Content, "/summarize") {
		go handleSummarizeCommand(s, m)
		return
	}

	// Handle command: /synthesizenow
	if strings.HasPrefix(m.Content, "/synthesizenow") {
		go handleSynthesizeNowCommand(s, m)
		return
	}

	// Handle command: /synthesize
	if strings.HasPrefix(m.Content, "/synthesize") {
		go handleSynthesizeCommand(s, m)
		return
	}

	// Handle command: /digest
	if strings.HasPrefix(m.Content, "/digest") {
		go handleDigestCommand(s, m)
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
	checkAndSummarizeMessages() // Consider scheduling this instead of calling immediately if needed
}

//// handleSummarizeCommand triggers manual summarization on command
//func handleSummarizeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
//	// Manually run the summarization process, ignoring the threshold
//	forceSummarizeMessages()
//
//	// Respond to the user
//	_, err := s.ChannelMessageSend(m.ChannelID, "Summarization triggered successfully, bypassing the threshold.")
//	if err != nil {
//		log.Printf("Error sending message: %v\n", err)
//	}
//}

// ---> attempt migrate refactor to summarize.go
//// forceSummarizeMessages summarizes all unsummarized messages, bypassing the threshold check
//func forceSummarizeMessages() {
//	// Fetch all unsummarized messages without checking for the threshold
//	rows, err := dbpool.Query(context.Background(),
//		`SELECT id, content FROM messages WHERE summarized = false`)
//	if err != nil {
//		log.Printf("Error fetching messages: %v\n", err)
//		return
//	}
//	defer rows.Close()
//
//	var contents []string
//	var messageIDs []int
//
//	// Accumulate all unsummarized messages
//	for rows.Next() {
//		var id int
//		var content string
//		err := rows.Scan(&id, &content)
//		if err != nil {
//			log.Printf("Error scanning message: %v\n", err)
//			continue
//		}
//
//		contents = append(contents, content)
//		messageIDs = append(messageIDs, id)
//	}
//
//	if len(contents) == 0 {
//		log.Println("No messages available for summarization.")
//		return
//	}
//
//	// Summarize the combined content using the OpenAI API
//	combinedContent := strings.Join(contents, "\n")
//	summary, err := summarizeContent(combinedContent)
//	if err != nil {
//		log.Printf("Error summarizing content: %v\n", err)
//		return
//	}
//
//	// Store the summary
//	_, err = dbpool.Exec(context.Background(),
//		`INSERT INTO summaries (summary) VALUES ($1)`, summary)
//	if err != nil {
//		log.Printf("Error inserting summary: %v\n", err)
//		return
//	}
//
//	// Mark the summarized messages
//	_, err = dbpool.Exec(context.Background(),
//		`UPDATE messages SET summarized = true WHERE id = ANY($1)`, messageIDs)
//	if err != nil {
//		log.Printf("Error updating messages: %v\n", err)
//	}
//
//	// Send the summary to the specified channel and tag everyone
//	err = sendSummaryToChannel(summarizeChannelID, summary)
//	if err != nil {
//		log.Printf("Error sending summary to channel: %v\n", err)
//	}
//
//	log.Printf("Successfully summarized all unsummarized messages without threshold.")
//}
//
//func checkAndSummarizeMessages() {
//	// Define threshold (e.g., 5000 characters)
//	const threshold = 5000
//
//	for {
//		// Fetch unsummarized messages
//		rows, err := dbpool.Query(context.Background(),
//			`SELECT id, content FROM messages WHERE summarized = false`)
//		if err != nil {
//			log.Printf("Error fetching messages: %v\n", err)
//			return
//		}
//
//		var contents []string
//		var messageIDs []int
//		var totalLength int
//
//		// Accumulate unsummarized messages
//		for rows.Next() {
//			var id int
//			var content string
//			err := rows.Scan(&id, &content)
//			if err != nil {
//				log.Printf("Error scanning message: %v\n", err)
//				continue
//			}
//
//			// Accumulate content length and details
//			totalLength += len(content)
//			contents = append(contents, content)
//			messageIDs = append(messageIDs, id)
//
//			// Break if the accumulated length reaches the threshold
//			if totalLength >= threshold {
//				break
//			}
//		}
//		rows.Close()
//
//		// Check if the length meets the threshold
//		if totalLength < threshold {
//			log.Printf("Accumulated content length (%d) is below the threshold (%d). No summarization needed.", totalLength, threshold)
//			return
//		}
//
//		// Summarize only when the threshold is met using the OpenAI API
//		combinedContent := strings.Join(contents, "\n")
//		summary, err := summarizeContent(combinedContent)
//		if err != nil {
//			log.Printf("Error summarizing content: %v\n", err)
//			return
//		}
//
//		// Store the summary
//		_, err = dbpool.Exec(context.Background(),
//			`INSERT INTO summaries (summary) VALUES ($1)`, summary)
//		if err != nil {
//			log.Printf("Error inserting summary: %v\n", err)
//			return
//		}
//
//		// Mark the summarized messages
//		_, err = dbpool.Exec(context.Background(),
//			`UPDATE messages SET summarized = true WHERE id = ANY($1)`, messageIDs)
//		if err != nil {
//			log.Printf("Error updating messages: %v\n", err)
//		}
//
//		// Send the summary to the specified channel and tag everyone
//		err = sendSummaryToChannel(summarizeChannelID, summary)
//		if err != nil {
//			log.Printf("Error sending summary to channel: %v\n", err)
//		}
//
//		log.Printf("Successfully summarized messages with accumulated length: %d.", totalLength)
//
//		// Reset accumulation variables after each summary
//		contents = nil
//		messageIDs = nil
//		totalLength = 0
//
//		// Exit the loop after summarizing
//		break
//	}
//}

// sendSummaryToChannel sends the summary to the specified channel and tags everyone
func sendSummaryToChannel(channelID, summary string) error {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}
	defer session.Close()

	// Send the summary message with @everyone tag
	_, err = session.ChannelMessageSend(channelID, fmt.Sprintf("@everyone\n**Summary:**\n\n%s", summary))
	if err != nil {
		return fmt.Errorf("error sending summary to channel %s: %w", channelID, err)
	}

	return nil
}
