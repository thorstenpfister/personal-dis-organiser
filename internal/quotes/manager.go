package quotes

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"personal-disorganizer/internal/parser"
)

// Manager handles quote loading and selection
type Manager struct {
	quotes    []parser.Quote
	configDir string
}

// NewManager creates a new quote manager
func NewManager(configDir string, quoteFiles []string) (*Manager, error) {
	m := &Manager{
		configDir: configDir,
		quotes:    []parser.Quote{},
	}
	
	// Load quotes from all configured files
	for _, file := range quoteFiles {
		if err := m.loadQuoteFile(file); err != nil {
			// Log error but continue - quotes are optional
			continue
		}
	}
	
	return m, nil
}

// loadQuoteFile loads quotes from a single file
func (m *Manager) loadQuoteFile(filename string) error {
	// Handle relative paths from config directory
	var filePath string
	if filepath.IsAbs(filename) {
		filePath = filename
	} else {
		filePath = filepath.Join(m.configDir, filename)
	}
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("quote file not found: %s", filePath)
	}
	
	// Load quotes
	quotes, err := parser.LoadQuotes(filePath)
	if err != nil {
		return fmt.Errorf("failed to load quotes from %s: %w", filePath, err)
	}
	
	m.quotes = append(m.quotes, quotes...)
	return nil
}

// GetRandomQuote returns a random quote
func (m *Manager) GetRandomQuote() *parser.Quote {
	if len(m.quotes) == 0 {
		return nil
	}
	
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(m.quotes))
	return &m.quotes[index]
}

// GetQuoteCount returns the total number of loaded quotes
func (m *Manager) GetQuoteCount() int {
	return len(m.quotes)
}

// HasQuotes returns true if quotes are available
func (m *Manager) HasQuotes() bool {
	return len(m.quotes) > 0
}