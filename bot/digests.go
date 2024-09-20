// bot/digest.go
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	openai "github.com/sashabaranov/go-openai"
)

// Replace with the actual ID of the channel where you want to send the digest
const digestChannelID = "1279545084412428410"

// generateDailyDigest creates a narrative summary of the last 24 hours, personalized for known users
func generateDailyDigest() {
	// Fetch undigested messages from the last 24 hours
	rows, err := dbpool.Query(context.Background(),
		`SELECT id, content, author_id FROM messages WHERE digested = false AND timestamp >= NOW() - INTERVAL '24 HOURS'`)
	if err != nil {
		log.Printf("Error fetching messages: %v\n", err)
		return
	}
	defer rows.Close()

	var contents []string
	var messageIDs []int

	for rows.Next() {
		var id int
		var content, authorID string
		err := rows.Scan(&id, &content, &authorID)
		if err != nil {
			log.Printf("Error scanning message: %v\n", err)
			continue
		}

		// Personalize the content by replacing author IDs with names
		content = personalizeContent(content, authorID)
		contents = append(contents, content)
		messageIDs = append(messageIDs, id)
	}

	if len(contents) == 0 {
		log.Println("No undigested messages available for digest generation.")
		return
	}

	// Combine messages into a single narrative
	combinedSummaries := strings.Join(contents, "\n")

	// Create a diary-style digest with personalizations
	digest, err := summarizeAsDiaryEntry(combinedSummaries)
	if err != nil {
		log.Printf("Error generating digest: %v\n", err)
		return
	}

	// Store the digest
	_, err = dbpool.Exec(context.Background(),
		`INSERT INTO digests (digest) VALUES ($1)`, digest)
	if err != nil {
		log.Printf("Error inserting digest: %v\n", err)
		return
	}

	// Mark the messages as digested
	_, err = dbpool.Exec(context.Background(),
		`UPDATE messages SET digested = true WHERE id = ANY($1)`, messageIDs)
	if err != nil {
		log.Printf("Error updating messages: %v\n", err)
	}

	// Send the digest to the specified channel
	err = sendMessageToChannel(digestChannelID, fmt.Sprintf("**Daily Digest:**\n\n%s", digest))
	if err != nil {
		log.Printf("Error sending digest to channel: %v\n", err)
	}

	log.Println("Daily digest generated and sent successfully!")
}

// summarizeAsDiaryEntry generates a third-person diary entry summarizing the "story so far" with personalized names
func summarizeAsDiaryEntry(content string) (string, error) {
	// Construct a prompt for a third-person diary-style summary
	prompt := fmt.Sprintf("You are writing a 3rd person narration summarizing recent events as if they are part of a character's journey in a video game. The characters are 2 lovers. Begin the entry with 'Story so far:' Stick to the facts provided without adding any fiction. Speculation and analysis can be allowed. \n\n%s", content)

	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini20240718, // or GPT-4 if available
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   1500, // Adjust based on expected digest length
			Temperature: 0.3,  // Set for structured, narrative output
		},
	)
	if err != nil {
		return "", err
	}

	// Return the diary-style narrative summary
	return resp.Choices[0].Message.Content, nil
}

// handleDigestCommand triggers manual digest generation on command
func handleDigestCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Manually run the digest generation process
	generateDailyDigest()

	// Respond to the user
	_, err := s.ChannelMessageSend(m.ChannelID, "Digest triggered successfully.")
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}
