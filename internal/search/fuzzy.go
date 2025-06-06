package search

import (
	"sort"
	"strings"
	"time"

	"personal-disorganizer/internal/storage"
)

// Result represents a search result
type Result struct {
	Task  storage.Task
	Score int
	Match string
}

// Engine handles fuzzy searching
type Engine struct{}

// NewEngine creates a new search engine
func NewEngine() *Engine {
	return &Engine{}
}

// Search performs fuzzy search across all tasks
func (e *Engine) Search(query string, tasks []storage.Task) []Result {
	if query == "" {
		return []Result{}
	}
	
	var results []Result
	query = strings.ToLower(query)
	today := time.Now().Truncate(24 * time.Hour)
	
	for _, task := range tasks {
		score := e.calculateScore(query, task.Text)
		if score > 0 {
			// Boost score for active/future tasks
			if !task.Done && !task.Date.Before(today) {
				score += 100
			}
			
			results = append(results, Result{
				Task:  task,
				Score: score,
				Match: e.highlightMatch(query, task.Text),
			})
		}
	}
	
	// Sort by score (highest first), then by date (newest first)
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score == results[j].Score {
			return results[i].Task.Date.After(results[j].Task.Date)
		}
		return results[i].Score > results[j].Score
	})
	
	return results
}

// calculateScore calculates fuzzy match score
func (e *Engine) calculateScore(query, text string) int {
	text = strings.ToLower(text)
	
	// Exact match gets highest score
	if strings.Contains(text, query) {
		if text == query {
			return 1000
		}
		return 500 + (100 - len(text)) // Prefer shorter matches
	}
	
	// Fuzzy matching - check if all characters in query appear in order
	queryChars := []rune(query)
	textChars := []rune(text)
	
	score := 0
	queryIdx := 0
	
	for i, char := range textChars {
		if queryIdx < len(queryChars) && char == queryChars[queryIdx] {
			// Characters match in order
			score += 10
			
			// Bonus for consecutive matches
			if queryIdx > 0 && i > 0 && textChars[i-1] == queryChars[queryIdx-1] {
				score += 5
			}
			
			// Bonus for word boundary matches
			if i == 0 || textChars[i-1] == ' ' {
				score += 15
			}
			
			queryIdx++
		}
	}
	
	// All query characters must be found
	if queryIdx < len(queryChars) {
		return 0
	}
	
	// Penalize for length difference
	score -= abs(len(textChars) - len(queryChars))
	
	return max(0, score)
}

// highlightMatch creates a highlighted version of the text
func (e *Engine) highlightMatch(query, text string) string {
	// Simple highlighting - just return the text for now
	// In a real implementation, you might add ANSI color codes
	return text
}

// abs returns absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// max returns maximum of two values
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}