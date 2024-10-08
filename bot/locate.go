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
