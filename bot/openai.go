// bot/openai.go
package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client

// initOpenAI initializes the OpenAI client with the API key
func initOpenAI() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	openaiClient = openai.NewClient(apiKey)
}

// -> attempt refactor to summarize.go
//// summarizeContent generates a reminder-focused summary of the provided content
//func summarizeContent(content string) (string, error) {
//	// Construct a reminder-focused prompt to generate accessible summaries
//	prompt := "You are a friendly assistant helping someone with Alzheimer's. Please create a gentle, easy-to-understand reminder based on the following information:\n\n" + content
//
//	resp, err := openaiClient.CreateChatCompletion(
//		context.Background(),
//		openai.ChatCompletionRequest{
//			Model: openai.GPT4oMini20240718, // Use GPT-3.5 Turbo or GPT-4 based on availability
//			Messages: []openai.ChatCompletionMessage{
//				{
//					Role:    openai.ChatMessageRoleUser,
//					Content: prompt,
//				},
//			},
//			MaxTokens:   1500, // Adjust based on expected summary length
//			Temperature: 0.7,  // Set temperature to balance friendliness and coherence
//		},
//	)
//	if err != nil {
//		return "", err
//	}
//
//	// Return the generated reminder-friendly summary
//	return resp.Choices[0].Message.Content, nil
//}

// personalizeContent replaces author IDs with names and fetches usernames for others
func personalizeContent(content, authorID string) string {
	// Replace known author IDs with personalized names
	switch authorID {
	case "869008800110243850":
		return "Hannah: " + content
	case "1123769580733603930":
		return "Jackson: " + content
	default:
		// Fetch the username from Discord for any other author ID
		username := fetchUsernameFromDiscord(authorID)
		return username + ": " + content
	}
}

// fetchUsernameFromDiscord fetches the Discord username given an author ID
func fetchUsernameFromDiscord(authorID string) string {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Printf("Error creating Discord session: %v\n", err)
		return "Unknown User"
	}
	defer session.Close()

	user, err := session.User(authorID)
	if err != nil {
		log.Printf("Error fetching user %s: %v\n", authorID, err)
		return "Unknown User"
	}

	return user.Username
}

// sendMessageToChannel sends a message to the specified channel using Discord session
func sendMessageToChannel(channelID, message string) error {
	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}
	defer session.Close()

	_, err = session.ChannelMessageSend(channelID, message)
	if err != nil {
		return fmt.Errorf("error sending message to channel %s: %w", channelID, err)
	}

	return nil
}
