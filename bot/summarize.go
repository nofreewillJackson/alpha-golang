// bot/summarize.go
package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"strings"
)

// handleSummarizeCommand triggers manual summarization on command
func handleSummarizeCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Manually run the summarization process, ignoring the threshold
	forceSummarizeMessages()

	// Respond to the user
	_, err := s.ChannelMessageSend(m.ChannelID, "Summarization triggered successfully, bypassing the threshold.")
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}

// summarizeContent generates a reminder-focused summary of the provided content
func summarizeContent(content string) (string, error) {
	// Construct a reminder-focused prompt to generate accessible summaries
	prompt := "You are a friendly assistant helping someone with Alzheimer's. Please create a gentle, easy-to-understand reminder based on the following information:\n\n" + content

	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini20240718, // Use GPT-3.5 Turbo or GPT-4 based on availability
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   1500, // Adjust based on expected summary length
			Temperature: 0.7,  // Set temperature to balance friendliness and coherence
		},
	)
	if err != nil {
		return "", err
	}

	// Return the generated reminder-friendly summary
	return resp.Choices[0].Message.Content, nil
}

// forceSummarizeMessages summarizes all unsummarized messages, bypassing the threshold check
func forceSummarizeMessages() {
	// Fetch all unsummarized messages without checking for the threshold
	rows, err := dbpool.Query(context.Background(),
		`SELECT id, content FROM messages WHERE summarized = false`)
	if err != nil {
		log.Printf("Error fetching messages: %v\n", err)
		return
	}
	defer rows.Close()

	var contents []string
	var messageIDs []int

	// Accumulate all unsummarized messages
	for rows.Next() {
		var id int
		var content string
		err := rows.Scan(&id, &content)
		if err != nil {
			log.Printf("Error scanning message: %v\n", err)
			continue
		}

		contents = append(contents, content)
		messageIDs = append(messageIDs, id)
	}

	if len(contents) == 0 {
		log.Println("No messages available for summarization.")
		return
	}

	// Summarize the combined content using the OpenAI API
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

	// Mark the summarized messages
	_, err = dbpool.Exec(context.Background(),
		`UPDATE messages SET summarized = true WHERE id = ANY($1)`, messageIDs)
	if err != nil {
		log.Printf("Error updating messages: %v\n", err)
	}

	// Send the summary to the specified channel and tag everyone
	err = sendSummaryToChannel(summarizeChannelID, summary)
	if err != nil {
		log.Printf("Error sending summary to channel: %v\n", err)
	}

	log.Printf("Successfully summarized all unsummarized messages without threshold.")
}

func checkAndSummarizeMessages() {
	// Define threshold (e.g., 5000 characters)
	const threshold = 5000

	for {
		// Fetch unsummarized messages
		rows, err := dbpool.Query(context.Background(),
			`SELECT id, content FROM messages WHERE summarized = false`)
		if err != nil {
			log.Printf("Error fetching messages: %v\n", err)
			return
		}

		var contents []string
		var messageIDs []int
		var totalLength int

		// Accumulate unsummarized messages
		for rows.Next() {
			var id int
			var content string
			err := rows.Scan(&id, &content)
			if err != nil {
				log.Printf("Error scanning message: %v\n", err)
				continue
			}

			// Accumulate content length and details
			totalLength += len(content)
			contents = append(contents, content)
			messageIDs = append(messageIDs, id)

			// Break if the accumulated length reaches the threshold
			if totalLength >= threshold {
				break
			}
		}
		rows.Close()

		// Check if the length meets the threshold
		if totalLength < threshold {
			log.Printf("Accumulated content length (%d) is below the threshold (%d). No summarization needed.", totalLength, threshold)
			return
		}

		// Summarize only when the threshold is met using the OpenAI API
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

		// Mark the summarized messages
		_, err = dbpool.Exec(context.Background(),
			`UPDATE messages SET summarized = true WHERE id = ANY($1)`, messageIDs)
		if err != nil {
			log.Printf("Error updating messages: %v\n", err)
		}

		// Send the summary to the specified channel and tag everyone
		err = sendSummaryToChannel(summarizeChannelID, summary)
		if err != nil {
			log.Printf("Error sending summary to channel: %v\n", err)
		}

		log.Printf("Successfully summarized messages with accumulated length: %d.", totalLength)

		// Reset accumulation variables after each summary
		contents = nil
		messageIDs = nil
		totalLength = 0

		// Exit the loop after summarizing
		break
	}
}
