// remind.go
package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"strconv"
	"strings"
)

func handleRemindCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore bot's own messages
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Trim the content after the `/remindme` command
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "/remindme"))

	// If the content is "show", display the reminders
	if content == "show" {
		displayReminders(s, m)
		return
	}

	// If the content is "clearall", delete all reminders
	if content == "clearall" {
		handleClearAllReminders(s, m)
		return
	}

	// If content is "delete ", handle deletion
	if strings.HasPrefix(content, "delete ") {
		handleDeleteReminder(s, m)
		return
	}

	// If the command is /remindme <text>, treat everything after /remindme as a reminder description
	if content != "" {
		addReminder(s, m, content)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please provide a reminder description.")
	}
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
	content := strings.TrimSpace(strings.TrimPrefix(m.Content, "/remindme delete "))

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

func sendHourlyReminders(dg *discordgo.Session) {
	var reminders []string
	var authorIDs []string

	rows, err := dbpool.Query(context.Background(),
		`SELECT description, author_id FROM reminders`)
	if err != nil {
		log.Printf("Error fetching reminders: %v\n", err)
		return
	}
	defer rows.Close()

	// Collect all reminders and author IDs
	for rows.Next() {
		var description, authorID string
		if err := rows.Scan(&description, &authorID); err != nil {
			log.Printf("Error scanning reminder: %v\n", err)
			continue
		}
		reminders = append(reminders, description)
		authorIDs = append(authorIDs, authorID)
	}

	// If there are no reminders, return early
	if len(reminders) == 0 {
		log.Println("No reminders to send.")
		return
	}

	// Prepare the message
	response := "@everyone Here are the current reminders:\n"
	for i, reminder := range reminders {
		response += fmt.Sprintf("%d. %s (by user %s)\n", i+1, reminder, authorIDs[i])
	}

	// Send to a specific channel (replace with your channel ID)
	channelID := os.Getenv("1293690503489126433") // Make sure to set this in your environment
	_, err = dg.ChannelMessageSend(channelID, response)
	if err != nil {
		log.Printf("Error sending reminder message: %v\n", err)
	} else {
		log.Println("Hourly reminders sent.")
	}
}
