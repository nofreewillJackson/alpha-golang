// bot/digest.go
package main

import (
	"context"
	"log"
	"strings"
)

func generateDailyDigest() {
	// Fetch summaries from the last 24 hours
	rows, err := dbpool.Query(context.Background(),
		`SELECT summary FROM summaries WHERE created_at >= NOW() - INTERVAL '24 HOURS'`)
	if err != nil {
		log.Printf("Error fetching summaries: %v\n", err)
		return
	}
	defer rows.Close()

	var summaries []string
	for rows.Next() {
		var summary string
		err := rows.Scan(&summary)
		if err != nil {
			log.Printf("Error scanning summary: %v\n", err)
			continue
		}
		summaries = append(summaries, summary)
	}

	if len(summaries) == 0 {
		return
	}

	// Summarize the summaries to create a digest
	combinedSummaries := strings.Join(summaries, "\n")
	digest, err := summarizeContent(combinedSummaries)
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

	log.Println("Daily digest generated successfully!")
}
