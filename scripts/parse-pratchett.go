package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Quote represents a single quote
type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	
	quotesDir := filepath.Join(homeDir, ".config", "personal-disorganizer", "quotes")
	inputFile := filepath.Join(quotesDir, "pratchett.pqf")
	outputFile := filepath.Join(quotesDir, "pratchett.json")
	
	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Input file not found: %s\n", inputFile)
		os.Exit(1)
	}
	
	quotes, err := parsePQF(inputFile)
	if err != nil {
		fmt.Printf("Error parsing PQF file: %v\n", err)
		os.Exit(1)
	}
	
	if err := saveQuotesToJSON(quotes, outputFile); err != nil {
		fmt.Printf("Error saving quotes to JSON: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("Successfully parsed %d quotes from %s to %s\n", len(quotes), inputFile, outputFile)
}

// parsePQF parses the Terry Pratchett quote file format
func parsePQF(filename string) ([]Quote, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	var quotes []Quote
	var currentQuote strings.Builder
	var currentAuthor string
	
	scanner := bufio.NewScanner(file)
	inQuote := false
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Empty line indicates end of quote
		if line == "" {
			if inQuote && currentQuote.Len() > 0 {
				quotes = append(quotes, Quote{
					Text:   strings.TrimSpace(currentQuote.String()),
					Author: currentAuthor,
				})
				currentQuote.Reset()
				currentAuthor = ""
				inQuote = false
			}
			continue
		}
		
		// Line starting with "-- " indicates attribution
		if strings.HasPrefix(line, "-- ") {
			currentAuthor = strings.TrimPrefix(line, "-- ")
			continue
		}
		
		// Regular quote text
		if currentQuote.Len() > 0 {
			currentQuote.WriteString(" ")
		}
		currentQuote.WriteString(line)
		inQuote = true
	}
	
	// Handle last quote if file doesn't end with empty line
	if inQuote && currentQuote.Len() > 0 {
		quotes = append(quotes, Quote{
			Text:   strings.TrimSpace(currentQuote.String()),
			Author: currentAuthor,
		})
	}
	
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	
	return quotes, nil
}

// saveQuotesToJSON saves quotes to a JSON file
func saveQuotesToJSON(quotes []Quote, filename string) error {
	data, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal quotes: %w", err)
	}
	
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}
	
	return nil
}