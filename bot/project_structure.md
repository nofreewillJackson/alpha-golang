# Project Directory Structure

```
bot/
├── bin/
    └── discordbot
├── .env
├── bot.exe
├── database.go
├── digests.go
├── go.mod
├── go.sum
├── handlers.go
├── llm.go
├── locate.go
├── log.go
├── main.go
├── openai.go
├── openai_test.go
├── printProject.py
├── project_structure.md
├── remind.go
├── summarize.go
├── synthesize.go
└── talk.go
```

# File Contents

## `C:\\Dev\\alpha-golang\\bot\\database.go`
```golang
// bot/database.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
)

var dbpool *pgxpool.Pool

func initDatabase() {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	var err error
	dbpool, err = pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Connected to the database successfully!")
}

```

## `C:\\Dev\\alpha-golang\\bot\\digests.go`
```golang
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

```

## `C:\\Dev\\alpha-golang\\bot\\handlers.go`
```golang
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

```

## `C:\\Dev\\alpha-golang\\bot\\llm.go`
```golang
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

```

## `C:\\Dev\\alpha-golang\\bot\\locate.go`
```golang
// locate.go
package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// Handle locate commands
func handleLocateCommands(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// /locate command
	if strings.HasPrefix(m.Content, "/locate ") {
		handleLocate(s, m)
		return
	}

	// /location command
	if strings.HasPrefix(m.Content, "/location ") {
		handleLocation(s, m)
		return
	}
}

// Handle /locate <item>
func handleLocate(s *discordgo.Session, m *discordgo.MessageCreate) {
	item := strings.TrimSpace(strings.TrimPrefix(m.Content, "/locate"))
	if item == "" {
		s.ChannelMessageSend(m.ChannelID, "Please specify the item you want to locate.")
		return
	}

	var description string
	err := dbpool.QueryRow(context.Background(),
		`SELECT description FROM locations WHERE item = $1`, item).Scan(&description)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Item not found in the locate database.")
		return
	}

	response := "Found in the locate database:\n\"" + description + "\""
	s.ChannelMessageSend(m.ChannelID, response)
}

// Handle /location <item>, <description>
func handleLocation(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "/location"))
	parts := strings.SplitN(content, ",", 2)
	if len(parts) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please use the format: `/location <item>, <description>`")
		return
	}

	item := strings.TrimSpace(parts[0])
	description := strings.TrimSpace(parts[1])

	_, err := dbpool.Exec(context.Background(),
		`INSERT INTO locations (item, description) VALUES ($1, $2)
         ON CONFLICT (item) DO UPDATE SET description = EXCLUDED.description`, item, description)
	if err != nil {
		log.Printf("Error inserting location: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, "Error adding location to the database.")
		return
	}

	response := "Confirmation: \"Location for '" + item + "' has been added successfully.\""
	s.ChannelMessageSend(m.ChannelID, response)
}

```

## `C:\\Dev\\alpha-golang\\bot\\log.go`
```golang
// log.go
package main

import (
	"context"
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

		response := "Your entry has been logged under " + formatTags(tags) + "."
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

// Categorize entry based on keywords
func categorizeEntry(entry string) []string {
	var tags []string
	lowerEntry := strings.ToLower(entry)

	if strings.Contains(lowerEntry, "today") || strings.Contains(lowerEntry, "feeling") {
		tags = append(tags, "#diary")
	}
	if strings.Contains(lowerEntry, "idea") || strings.Contains(lowerEntry, "brainstormed") {
		tags = append(tags, "#ideas")
	}
	if strings.Contains(lowerEntry, "need to") || strings.Contains(lowerEntry, "finish") {
		tags = append(tags, "#tasks")
	}
	if strings.Contains(lowerEntry, "anxious") || strings.Contains(lowerEntry, "stress") || strings.Contains(lowerEntry, "feeling") {
		tags = append(tags, "#emotions")
	}
	if len(tags) == 0 {
		tags = append(tags, "#general")
	}
	return tags
}

func formatTags(tags []string) string {
	return strings.Join(tags, ", ")
}

```

## `C:\\Dev\\alpha-golang\\bot\\main.go`
```golang
// bot/main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	// Load environment variables from .env file
	//err := godotenv.Load()
	//if err != nil {
	//	log.Println("No .env file found, relying on environment variables")
	//} else {
	//	log.Println(".env file loaded successfully")
	//}
	err := godotenv.Overload() // Overload the .env if in system file
	if err != nil {
		log.Printf("Error overloading .env file: %v", err)
	} else {
		log.Println(".env file loaded and variables overwritten")
	}
	// Initialize database connection
	initDatabase()
	defer dbpool.Close()

	// Initialize OpenAI client
	initOpenAI()

	// Create a new Discord session using the provided bot token
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_BOT_TOKEN is not set. Please check your .env file or environment variables.")
	} else {
		log.Printf("Discord Bot Token is set. Token length: %d characters", len(token))
		// Be careful not to log the actual token for security reasons
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// Register the messageCreate function as a callback for MessageCreate events
	dg.AddHandler(messageCreate)
	dg.AddHandler(handleLocateCommands)
	dg.AddHandler(handleLogCommand)
	dg.AddHandler(handleRemindCommand)
	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	// Schedule daily digest generation
	c := cron.New()
	_, err = c.AddFunc("@daily", generateDailyDigest)
	if err != nil {
		log.Fatalf("Error scheduling daily digest: %v", err)
	}
	c.Start()
	defer c.Stop()

	log.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session
	dg.Close()
}

```

## `C:\\Dev\\alpha-golang\\bot\\openai.go`
```golang
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

```

## `C:\\Dev\\alpha-golang\\bot\\openai_test.go`
```golang
// bot/openai_test.go
package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestSummarizeContent tests the summarizeContent function
func TestSummarizeContent(t *testing.T) {
	// Initialize the OpenAI client
	initOpenAI()

	// Define test content to summarize
	content := "how rare is a black rabbit?"

	// Call the summarizeContent function
	summary, err := summarizeContent(content)

	// Check if there was an error
	if err != nil {
		t.Fatalf("Error summarizing content: %v", err)
	}

	// Print the returned summary for inspection
	fmt.Printf("Summary: %s\n", summary)

	// Check if the summary is empty
	if summary == "" {
		t.Error("Summary should not be empty")
	}

	// Check if the summary length is reasonable (example: less than 100 characters)
	if len(summary) > 100 {
		t.Errorf("Summary is unexpectedly long (%d characters): %s", len(summary), summary)
	}

	// Additional checks (optional):
	// You can add more specific checks depending on the expected behavior
	// For example, checking if certain keywords from the original content appear in the summary.
	if !containsExpectedKeywords(summary, []string{"test", "summarize"}) {
		t.Error("Summary does not contain expected keywords.")
	}
}

// Helper function to check if the summary contains expected keywords
func containsExpectedKeywords(summary string, keywords []string) bool {
	for _, keyword := range keywords {
		if !containsIgnoreCase(summary, keyword) {
			return false
		}
	}
	return true
}

// Helper function to check case-insensitive containment
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

```

## `C:\\Dev\\alpha-golang\\bot\\remind.go`
```golang
// remind.go
package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
)

func handleRemindCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "/remind" {
		displayReminders(s, m)
		return
	}

	// Handle 'delete' and 'later' commands
	if strings.HasPrefix(m.Content, "delete ") {
		handleDeleteReminder(s, m)
		return
	}
	if strings.HasPrefix(m.Content, "later ") {
		handleLaterReminder(s, m)
		return
	}
}

// Display current reminders
func displayReminders(s *discordgo.Session, m *discordgo.MessageCreate) {
	rows, err := dbpool.Query(context.Background(),
		`SELECT id, description FROM reminders WHERE author_id = $1 ORDER BY id`, m.Author.ID)
	if err != nil {
		log.Printf("Error fetching reminders: %v\n", err)
		return
	}
	defer rows.Close()

	var reminders []string
	for rows.Next() {
		var id int
		var description string
		err := rows.Scan(&id, &description)
		if err != nil {
			log.Printf("Error scanning reminder: %v\n", err)
			continue
		}
		reminders = append(reminders, strconv.Itoa(id)+". "+description)
	}

	if len(reminders) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You have no current reminders.")
		return
	}

	response := "Reminders:\n" + strings.Join(reminders, "\n") + "\n\nOptions:\n- Type 'delete <number>' to remove a reminder.\n- Type 'later <number>' to postpone a reminder."
	s.ChannelMessageSend(m.ChannelID, response)
}

// Handle 'delete' command
func handleDeleteReminder(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "delete"))
	id, err := strconv.Atoi(content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please provide a valid reminder number to delete.")
		return
	}

	result, err := dbpool.Exec(context.Background(),
		`DELETE FROM reminders WHERE id = $1 AND author_id = $2`, id, m.Author.ID)
	if err != nil {
		log.Printf("Error deleting reminder: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, "Error deleting the reminder.")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		s.ChannelMessageSend(m.ChannelID, "No reminder found with that number.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Confirmation: Reminder "+strconv.Itoa(id)+" has been deleted.")
}

// Handle 'later' command
func handleLaterReminder(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "later"))
	id, err := strconv.Atoi(content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please provide a valid reminder number to postpone.")
		return
	}

	// Update the reminder's due date (assuming a 'due_date' column exists)
	_, err = dbpool.Exec(context.Background(),
		`UPDATE reminders SET due_date = NOW() + INTERVAL '1 DAY' WHERE id = $1 AND author_id = $2`, id, m.Author.ID)
	if err != nil {
		log.Printf("Error postponing reminder: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, "Error postponing the reminder.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Confirmation: Reminder "+strconv.Itoa(id)+" has been postponed to tomorrow.")
}

```

## `C:\\Dev\\alpha-golang\\bot\\summarize.go`
```golang
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
	prompt := "You are a friendly assistant helping someone with Alzheimer's(but you MUST BE IN THIRD PERSON). Please create a gentle, easy-to-understand reminder based on the following information:\n\n" + content

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

```

## `C:\\Dev\\alpha-golang\\bot\\synthesize.go`
```golang
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
		"You are a couples' counselor. Your job is to digest the following messages and help %s communicate better with %s. Fit yourself in %s's shoes(but you MUST BE IN THIRD PERSON), make their grievances and perspective coherent and understandable to %s. Provide empathetic and compassionate feedback to help resolve their concerns.\n\n%s",
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

```

## `C:\\Dev\\alpha-golang\\bot\\talk.go`
```golang
package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// Replace with the actual ID of the channel where the bot will "talk"
const talkChannelID = "1286554426768883753"

// handleTalkCommand handles the /talk command and makes the bot send the message to the channel
func handleTalkCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Extract the content after the /talk command
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "/talk"))

	// If no content is provided, send an error message
	if content == "" {
		_, err := s.ChannelMessageSend(m.ChannelID, "You need to provide something for me to say!")
		if err != nil {
			log.Printf("Error sending message: %v\n", err)
		}
		return
	}

	// Send the message to the designated channel
	_, err := s.ChannelMessageSend(talkChannelID, content)
	if err != nil {
		log.Printf("Error sending message: %v\n", err)
		return
	}

	// Optionally, acknowledge the command to the user
	_, err = s.ChannelMessageSend(m.ChannelID, "Message sent!")
	if err != nil {
		log.Printf("Error sending acknowledgment: %v\n", err)
	}
}

```

