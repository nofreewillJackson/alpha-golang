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

	// Open a websocket connection to Discord and begin listening
	err = dg.Open()
	if err != nil {
		log.Fatalf("Error opening Discord session: %v", err)
	}

	// Initialize the cron scheduler
	c := cron.New()

	// Schedule hourly reminder ping
	c.AddFunc("@hourly", func() {
		sendHourlyReminders(dg)
	})

	// Schedule daily digest generation
	_, err = c.AddFunc("@daily", generateDailyDigest)
	if err != nil {
		log.Fatalf("Error scheduling daily digest: %v", err)
	}

	// Start the cron jobs
	c.Start()
	defer c.Stop()

	log.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session
	dg.Close()
}
