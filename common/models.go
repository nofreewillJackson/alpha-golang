// common/models/models.go
package common

import "time"

// Message represents a Discord message stored in the database.
type Message struct {
	ID         int       `json:"id"`
	Content    string    `json:"content"`
	AuthorID   string    `json:"author_id"`
	ChannelID  string    `json:"channel_id"`
	Timestamp  time.Time `json:"timestamp"`
	Summarized bool      `json:"summarized"`
}

// Summary represents a summarized text stored in the database.
type Summary struct {
	ID        int       `json:"id"`
	Summary   string    `json:"summary"`
	CreatedAt time.Time `json:"created_at"`
}

// Digest represents a daily digest stored in the database.
type Digest struct {
	ID        int       `json:"id"`
	Digest    string    `json:"digest"`
	CreatedAt time.Time `json:"created_at"`
}
