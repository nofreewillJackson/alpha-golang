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

	// Handle /nuke command
	if strings.HasPrefix(m.Content, "/nuke") {
		go handleNukeCommand(s, m)
		return
	}

	// Handle /locate command
	if strings.HasPrefix(m.Content, "/locate ") {
		go handleLocate(s, m)
		return
	}

	// Handle /location command
	if strings.HasPrefix(m.Content, "/location ") {
		go handleLocation(s, m)
		return
	}

	// Handle /remind command
	if strings.HasPrefix(m.Content, "/remindme") {
		go handleRemindCommand(s, m)
		return
	}

	if strings.HasPrefix(m.Content, "/log") {
		go handleLogCommand(s, m)
		return
	}

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
func handleNukeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Let the user know the nuke is starting
	s.ChannelMessageSend(m.ChannelID, "Nuking the channel... 💣")

	// Loop through batches of messages
	for {
		// Fetch up to 100 messages at a time
		messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		if err != nil {
			log.Printf("Error fetching messages: %v\n", err)
			s.ChannelMessageSend(m.ChannelID, "Error fetching messages.")
			return
		}

		// Break if there are no messages left
		if len(messages) == 0 {
			break
		}

		// Iterate through the messages to delete them one by one if they are older than 2 weeks
		for _, msg := range messages {
			// Use `ChannelMessageDelete` for messages that can't be bulk deleted (older than 2 weeks)
			err := s.ChannelMessageDelete(m.ChannelID, msg.ID)
			if err != nil {
				log.Printf("Error deleting message: %v\n", err)
				continue
			}

			// Be mindful of Discord rate limits
			// Adding a short delay to avoid hitting the rate limit
			time.Sleep(500 * time.Millisecond)
		}
	}

	// Send confirmation message after nuking
	s.ChannelMessageSend(m.ChannelID, "Channel nuked successfully!")
}

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
