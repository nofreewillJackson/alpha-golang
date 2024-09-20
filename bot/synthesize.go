// bot/synthesize.go
package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// Replace with the actual ID of the channel where synthesized messages will be sent
const synthesizeChannelID = "1286554426768883753"

// handleSynthesizeNowCommand triggers synthesis based on the user's input following the /synthesizenow command
func handleSynthesizeNowCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Extract the content after the /synthesizenow command
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "/synthesizenow"))

	// If no content is provided, send an error message
	if content == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "Jackson: you're supposed to provide context after the prefix, dummy.")
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
		return
	}

	// Determine the sender's perspective based on their user ID
	var recipient string
	switch m.Author.ID {
	case "1123769580733603930": // Jackson's ID
		recipient = "Hannah"
	case "869008800110243850": // Hannah's ID
		recipient = "Jackson"
	}

	// Generate a synthesis response using the provided content
	synthesis, err := synthesizeMessagesAsCouplesCounselor(content, m.Author.Username, recipient)
	if err != nil {
		log.Printf("Error generating synthesis: %v\n", err)
		_, err := s.ChannelMessageSend(m.ChannelID, "ask jackson. error synthesizing")
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
		return
	}

	// Send the synthesized response back to the Discord channel
	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**Synthesis:**\n\n%s", synthesis))
	if err != nil {
		log.Printf("Error sending synthesized message: %v\n", err)
	}
}

// synthesizeMessagesAsCouplesCounselor generates a relationship counselor-style synthesis of the messages
func synthesizeMessagesAsCouplesCounselor(content, sender, recipient string) (string, error) {
	// Construct a prompt to simulate a couples' counselor providing insights
	prompt := fmt.Sprintf(
		"You are a couples' counselor. Your job is to digest the following messages and help %s communicate better with %s. Fit yourself in %s's shoes, make their grievances and perspective coherent and understandable to %s. Provide empathetic and compassionate feedback to help resolve their concerns.\n\n%s",
		sender, recipient, sender, recipient, content,
	)

	// Use the OpenAI client to generate the synthesis
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
			MaxTokens:   1500, // Adjust based on expected synthesis length
			Temperature: 0.7,  // Adjust for empathetic, nuanced output
		},
	)
	if err != nil {
		return "", err
	}

	// Return the counselor-style synthesis
	return resp.Choices[0].Message.Content, nil
}

// handleSynthesizeCommand triggers manual synthesis generation based on the last 24 hours of messages
func handleSynthesizeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Manually run the synthesis generation process
	generateSynthesis()

	// Respond to the user
	_, err := s.ChannelMessageSend(m.ChannelID, "*trying my best* ~ jackson temp mod v.0.0.8")
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}

// generateSynthesis creates a relationship counselor-style synthesis of the recent messages between two people
func generateSynthesis() {
	// Fetch unsynthesized messages from the specific synthesize channel in the database
	rows, err := dbpool.Query(context.Background(),
		`SELECT id, content, author_id FROM messages WHERE synthesized = false AND channel_id = $1 AND timestamp >= NOW() - INTERVAL '24 HOURS'`, synthesizeChannelID)
	if err != nil {
		log.Printf("Error fetching messages: %v\n", err)
		return
	}
	defer rows.Close()

	var contents []string
	var messageIDs []int
	var sender, recipient string

	for rows.Next() {
		var id int
		var content, authorID string
		err := rows.Scan(&id, &content, &authorID)
		if err != nil {
			log.Printf("Error scanning message: %v\n", err)
			continue
		}

		// Determine sender and recipient based on authorID
		switch authorID {
		case "1123769580733603930": // Jackson's ID
			sender = "Jackson"
			recipient = "Hannah"
		case "869008800110243850": // Hannah's ID
			sender = "Hannah"
			recipient = "Jackson"
		}

		// Personalize the content by replacing author IDs with names
		content = personalizeContent(content, authorID)
		contents = append(contents, content)
		messageIDs = append(messageIDs, id)
	}

	if len(contents) == 0 {
		log.Println("No unsynthesized messages available for synthesis generation.")
		return
	}

	// Combine messages into a single narrative
	combinedMessages := strings.Join(contents, "\n")

	// Create a relationship-counselor-style synthesis
	synthesis, err := synthesizeMessagesAsCouplesCounselor(combinedMessages, sender, recipient)
	if err != nil {
		log.Printf("Error generating synthesis: %v\n", err)
		return
	}

	// Store the synthesis in the messages table and mark as synthesized
	_, err = dbpool.Exec(context.Background(),
		`UPDATE messages SET synthesis = $1, synthesized = true WHERE id = ANY($2)`, synthesis, messageIDs)
	if err != nil {
		log.Printf("Error updating messages with synthesis: %v\n", err)
	}

	// Send the synthesis to the specified channel
	err = sendMessageToChannel(synthesizeChannelID, fmt.Sprintf("**Synthesis:**\n\n%s", synthesis))
	if err != nil {
		log.Printf("Error sending synthesis to channel: %v\n", err)
	}

	log.Println("Synthesis generated and sent successfully!")
}
