// bot/openai.go
package main

import (
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client

// initOpenAI initializes the OpenAI client with the API key
func initOpenAI() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	openaiClient = openai.NewClient(apiKey)
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
