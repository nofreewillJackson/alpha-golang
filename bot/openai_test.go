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
