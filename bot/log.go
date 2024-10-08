// log.go
package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func handleLogCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "/log ") {
		content := strings.TrimSpace(strings.TrimPrefix(m.Content, "/log"))
		if content == "" {
			s.ChannelMessageSend(m.ChannelID, "Please provide an entry to log.")
			return
		}

		// Automatic categorization
		tags := categorizeEntry(content)

		_, err := dbpool.Exec(context.Background(),
			`INSERT INTO logs (entry, tags, author_id) VALUES ($1, $2, $3)`, content, strings.Join(tags, ", "), m.Author.ID)
		if err != nil {
			log.Printf("Error inserting log entry: %v\n", err)
			s.ChannelMessageSend(m.ChannelID, "Error logging the entry.")
			return
		}

		response := fmt.Sprintf("Your entry has been logged under %s.", formatTags(tags))
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

// Categorize entry based on keywords
func categorizeEntry(entry string) []string {
	var tags []string
	lowerEntry := strings.ToLower(entry)

	if strings.Contains(lowerEntry, "today") || strings.Contains(lowerEntry, "yesterday") {
		tags = append(tags, "#diary")
	}
	if strings.Contains(lowerEntry, "idea") || strings.Contains(lowerEntry, "brainstormed") {
		tags = append(tags, "#ideas")
	}
	if strings.Contains(lowerEntry, "need to") || strings.Contains(lowerEntry, "finish") || strings.Contains(lowerEntry, "must") {
		tags = append(tags, "#tasks")
	}
	if strings.Contains(lowerEntry, "anxious") || strings.Contains(lowerEntry, "stress") || strings.Contains(lowerEntry, "feeling") || strings.Contains(lowerEntry, "happy") || strings.Contains(lowerEntry, "sad") {
		tags = append(tags, "#emotions")
	}
	if strings.Contains(lowerEntry, "project") {
		tags = append(tags, "#project")
	}
	if len(tags) == 0 {
		tags = append(tags, "#general")
	}
	return tags
}

func formatTags(tags []string) string {
	return strings.Join(tags, ", ")
}
