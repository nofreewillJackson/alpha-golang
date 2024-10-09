// remind.go
package main

import (
	"context"
	"fmt"
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

	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "/reminders"))
	if content == "" {
		displayReminders(s, m)
		return
	}

	if strings.HasPrefix(content, "delete ") {
		handleDeleteReminder(s, m)
		return
	}

	if strings.HasPrefix(content, "later ") {
		handleLaterReminder(s, m)
		return
	}

	if strings.HasPrefix(content, "clearall") {
		handleClearAllReminders(s, m)
		return
	}

	// Otherwise, add a new reminder
	addReminder(s, m, content)
}

// Add a new reminder
func addReminder(s *discordgo.Session, m *discordgo.MessageCreate, description string) {
	_, err := dbpool.Exec(context.Background(),
		`INSERT INTO reminders (description, author_id) VALUES ($1, $2)`, description, m.Author.ID)
	if err != nil {
		log.Printf("Error adding reminder: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, "Error adding the reminder.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Your reminder has been added.")
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
		reminders = append(reminders, fmt.Sprintf("%d. %s", id, description))
	}

	if len(reminders) == 0 {
		s.ChannelMessageSend(m.ChannelID, "You have no current reminders.")
		return
	}

	response := "Reminders:\n" + strings.Join(reminders, "\n") + "\n\nOptions:\n" +
		"- Type 'delete <number>' to remove a reminder.\n" +
		"- Type 'later <number>' to postpone a reminder."

	s.ChannelMessageSend(m.ChannelID, response)
}

func handleDeleteReminder(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ensure the message starts with "delete " and extract the number after it
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "/remind delete "))

	// Convert the string to an integer (the reminder number to delete)
	id, err := strconv.Atoi(content)
	if err != nil || id <= 0 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a valid reminder number to delete.")
		return
	}

	// Execute the delete operation and check if a row was affected
	result, err := dbpool.Exec(context.Background(),
		`DELETE FROM reminders WHERE id = $1 AND author_id = $2`, id, m.Author.ID)
	if err != nil {
		log.Printf("Error deleting reminder: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, "Error deleting the reminder.")
		return
	}

	// Check if any rows were affected
	if result.RowsAffected() == 0 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No reminder found with number %d.", id))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Confirmation: Reminder %d has been deleted.", id))
}

// handleClearAllReminders deletes all reminders for the user or all reminders in the system
func handleClearAllReminders(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Clear only the user's reminders (optional, you can modify this to clear all reminders if desired)
	result, err := dbpool.Exec(context.Background(),
		`DELETE FROM reminders WHERE author_id = $1`, m.Author.ID)
	if err != nil {
		log.Printf("Error clearing reminders: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, "Error clearing reminders.")
		return
	}

	// Check if any rows were affected (i.e., reminders were deleted)
	if result.RowsAffected() == 0 {
		s.ChannelMessageSend(m.ChannelID, "No reminders found to clear.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "All your reminders have been cleared.")
}

// Handle 'later' command
func handleLaterReminder(s *discordgo.Session, m *discordgo.MessageCreate) {
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "later "))
	id, err := strconv.Atoi(content)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Please provide a valid reminder number to postpone.")
		return
	}

	// Update the reminder's due date (assuming a 'due_date' column exists)
	_, err = dbpool.Exec(context.Background(),
		`UPDATE reminders SET due_date = due_date + INTERVAL '1 DAY' WHERE id = $1 AND author_id = $2`, id, m.Author.ID)
	if err != nil {
		log.Printf("Error postponing reminder: %v\n", err)
		s.ChannelMessageSend(m.ChannelID, "Error postponing the reminder.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Confirmation: Reminder %d has been postponed to tomorrow.", id))
}
