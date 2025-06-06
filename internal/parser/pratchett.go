package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Quote represents a single quote
type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

// ParsePQF parses the Terry Pratchett quote file format from a file
func ParsePQF(filename string) ([]Quote, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	return parsePQFReader(file)
}

// parsePQFReader parses the PQF format from any reader
func parsePQFReader(file *os.File) ([]Quote, error) {
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

// LoadQuotes loads quotes from a JSON file
func LoadQuotes(filename string) ([]Quote, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read quotes file: %w", err)
	}
	
	var quotes []Quote
	if err := json.Unmarshal(data, &quotes); err != nil {
		return nil, fmt.Errorf("failed to parse quotes JSON: %w", err)
	}
	
	return quotes, nil
}