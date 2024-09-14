// bot/openai_test.go
package main

import (
	"testing"
)

func TestSummarizeContent(t *testing.T) {
	initOpenAI()
	content := "This is a test content to summarize."
	summary, err := summarizeContent(content)
	if err != nil {
		t.Errorf("Error summarizing content: %v", err)
	}
	if summary == "" {
		t.Error("Summary should not be empty")
	}
}
