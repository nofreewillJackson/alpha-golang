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
