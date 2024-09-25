package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// Replace with the actual ID of the channel where the LLM response will be sent
const llmChannelID = "1286554426768883753"

// handleLLMCommand triggers ChatGPT based on the user's input following the /llm command
func handleLLMCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Extract the content after the /llm command
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "/llm"))

	// If no content is provided, send an error message
	if content == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "You're supposed to provide content after the /llm command.")
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
		return
	}

	// Generate a response using the provided content
	llmResponse, err := queryOpenAI(content)
	if err != nil {
		log.Printf("Error generating LLM response: %v\n", err)
		_, err := s.ChannelMessageSend(m.ChannelID, "Error processing the LLM request.")
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
		return
	}

	// Send the generated response back to the designated channel
	_, err = s.ChannelMessageSend(llmChannelID, fmt.Sprintf("**LLM Response:**\n\n%s", llmResponse))
	if err != nil {
		log.Printf("Error sending LLM response: %v\n", err)
	}
}

// queryOpenAI sends the provided content to OpenAI's API and returns the response
func queryOpenAI(content string) (string, error) {
	// Use the OpenAI client to generate the response
	resp, err := openaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4, // or GPT-4 if available
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
			MaxTokens:   1500, // Adjust based on expected response length
			Temperature: 0.7,  // Adjust for nuanced and creative output
		},
	)
	if err != nil {
		return "", err
	}

	// Return the generated LLM response
	return resp.Choices[0].Message.Content, nil
}
