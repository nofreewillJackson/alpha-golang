package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/http"
	"strings"
)

type NLPResponse struct {
	Sentiment float64  `json:"sentiment"`
	Entities  []string `json:"entities"` // Dynamic entity tags from Python service
}

func callNLPService(text string) (*NLPResponse, error) {
	payload := map[string]string{"text": text}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post("http://localhost:5000/analyze", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result NLPResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

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

		// Store tags as JSONB
		tagsJSON, err := json.Marshal(tags)
		if err != nil {
			log.Printf("Error marshaling tags to JSON: %v\n", err)
			s.ChannelMessageSend(m.ChannelID, "Error processing tags.")
			return
		}

		_, err = dbpool.Exec(context.Background(),
			`INSERT INTO logs (entry, tags, author_id) VALUES ($1, $2::jsonb, $3)`, content, tagsJSON, m.Author.ID)
		if err != nil {
			log.Printf("Error inserting log entry: %v\n", err)
			s.ChannelMessageSend(m.ChannelID, "Error logging the entry.")
			return
		}

		response := fmt.Sprintf("Your entry has been logged under %s.", formatTags(tags))
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

// Categorize entry using NLP and keyword-based logic
func categorizeEntry(entry string) []string {
	var tags []string
	lowerEntry := strings.ToLower(entry)

	// Call NLP Service for dynamic tagging
	nlpResponse, err := callNLPService(entry)
	if err != nil {
		log.Printf("Error calling NLP service: %v\n", err)
		tags = append(tags, "#general")
	} else {
		// Sentiment Analysis Tags
		if nlpResponse.Sentiment > 0.5 {
			tags = append(tags, "#emotions/positive")
		} else if nlpResponse.Sentiment < -0.5 {
			tags = append(tags, "#emotions/negative")
		} else {
			tags = append(tags, "#emotions/neutral")
		}

		// Add dynamic entity tags from the NLP service
		tags = append(tags, nlpResponse.Entities...)
	}

	// Existing Keyword-Based Tags
	if strings.Contains(lowerEntry, "today") || strings.Contains(lowerEntry, "yesterday") {
		tags = append(tags, "#diary")
	}
	if strings.Contains(lowerEntry, "idea") || strings.Contains(lowerEntry, "brainstormed") {
		tags = append(tags, "#ideas")
	}
	if strings.Contains(lowerEntry, "need to") || strings.Contains(lowerEntry, "finish") || strings.Contains(lowerEntry, "must") {
		tags = append(tags, "#tasks/urgent")
	} else if strings.Contains(lowerEntry, "should") || strings.Contains(lowerEntry, "might") {
		tags = append(tags, "#tasks/non-urgent")
	}
	if strings.Contains(lowerEntry, "project") {
		tags = append(tags, "#project")
	}
	if len(tags) == 0 {
		tags = append(tags, "#general")
	}

	return tags
}

// Format tags for display
func formatTags(tags []string) string {
	return strings.Join(tags, ", ")
}
