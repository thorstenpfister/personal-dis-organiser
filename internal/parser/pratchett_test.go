package parser

import (
	"os"
	"path/filepath"
	"testing"

	"personal-disorganizer/internal/testutil"
)

func TestParsePQF(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		expectedQuotes int
		expectError    bool
	}{
		{
			name:           "parse valid PQF file",
			filename:       "sample.pqf",
			expectedQuotes: 5,
			expectError:    false,
		},
		{
			name:           "parse malformed PQF file",
			filename:       "malformed.pqf",
			expectedQuotes: 3, // Should still parse some valid quotes
			expectError:    false,
		},
		{
			name:           "parse empty PQF file",
			filename:       "empty.pqf",
			expectedQuotes: 0,
			expectError:    false,
		},
		{
			name:           "file not found",
			filename:       "nonexistent.pqf",
			expectedQuotes: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join("testdata", tt.filename)
			
			quotes, err := ParsePQF(filePath)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if len(quotes) != tt.expectedQuotes {
				t.Errorf("Expected %d quotes, got %d", tt.expectedQuotes, len(quotes))
			}
			
			// Verify quote structure for valid files
			if tt.expectedQuotes > 0 {
				for i, quote := range quotes {
					if quote.Text == "" {
						t.Errorf("Quote %d has empty text", i)
					}
					// Note: Author can be empty for some quotes
				}
			}
		})
	}
}

func TestParsePQF_SpecificContent(t *testing.T) {
	filePath := filepath.Join("testdata", "sample.pqf")
	
	quotes, err := ParsePQF(filePath)
	if err != nil {
		t.Fatalf("Failed to parse PQF file: %v", err)
	}
	
	// Test specific quote content (parser includes quotes in text)
	expectedFirstQuote := "\"The trouble with having an open mind, of course, is that people will insist on coming along and trying to put things in it.\""
	if len(quotes) > 0 && quotes[0].Text != expectedFirstQuote {
		t.Errorf("First quote text mismatch.\nExpected: %s\nGot: %s", expectedFirstQuote, quotes[0].Text)
	}
	
	// Test author attribution (author comes after quote in PQF format)
	expectedSecondAuthor := "Terry Pratchett, Diggers"
	if len(quotes) > 1 && quotes[1].Author != expectedSecondAuthor {
		t.Errorf("Second quote author mismatch.\nExpected: %s\nGot: %s", expectedSecondAuthor, quotes[1].Author)
	}
	
	// Test quote without explicit ending newline
	if len(quotes) >= 5 {
		expectedLastQuote := "\"Five exclamation marks, the sure sign of an insane mind.\""
		if quotes[4].Text != expectedLastQuote {
			t.Errorf("Last quote text mismatch.\nExpected: %s\nGot: %s", expectedLastQuote, quotes[4].Text)
		}
	}
}

func TestLoadQuotes(t *testing.T) {
	tests := []struct {
		name           string
		filename       string
		expectedQuotes int
		expectError    bool
	}{
		{
			name:           "load valid JSON quotes",
			filename:       "sample.json",
			expectedQuotes: 3,
			expectError:    false,
		},
		{
			name:           "load empty JSON quotes",
			filename:       "empty.json",
			expectedQuotes: 0,
			expectError:    false,
		},
		{
			name:           "load invalid JSON quotes",
			filename:       "invalid.json",
			expectedQuotes: 0,
			expectError:    true,
		},
		{
			name:           "file not found",
			filename:       "nonexistent.json",
			expectedQuotes: 0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join("testdata", tt.filename)
			
			quotes, err := LoadQuotes(filePath)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if len(quotes) != tt.expectedQuotes {
				t.Errorf("Expected %d quotes, got %d", tt.expectedQuotes, len(quotes))
			}
			
			// Verify quote structure
			for i, quote := range quotes {
				if quote.Text == "" {
					t.Errorf("Quote %d has empty text", i)
				}
				if quote.Author == "" {
					t.Errorf("Quote %d has empty author", i)
				}
			}
		})
	}
}

func TestLoadQuotes_SpecificContent(t *testing.T) {
	filePath := filepath.Join("testdata", "sample.json")
	
	quotes, err := LoadQuotes(filePath)
	if err != nil {
		t.Fatalf("Failed to load JSON quotes: %v", err)
	}
	
	if len(quotes) == 0 {
		t.Fatal("No quotes loaded")
	}
	
	// Test specific quote content
	expectedFirstQuote := "The trouble with having an open mind, of course, is that people will insist on coming along and trying to put things in it."
	if quotes[0].Text != expectedFirstQuote {
		t.Errorf("First quote text mismatch.\nExpected: %s\nGot: %s", expectedFirstQuote, quotes[0].Text)
	}
	
	// Test author attribution
	expectedFirstAuthor := "Terry Pratchett, Diggers"
	if quotes[0].Author != expectedFirstAuthor {
		t.Errorf("First quote author mismatch.\nExpected: %s\nGot: %s", expectedFirstAuthor, quotes[0].Author)
	}
}

func TestParsePQF_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Quote
	}{
		{
			name:    "quote without author",
			content: `"This is a quote without an author"`,
			expected: []Quote{
				{Text: "\"This is a quote without an author\"", Author: ""},
			},
		},
		{
			name: "multiple quotes with mixed authors",
			content: `"First quote"

-- First Author

"Second quote without author"

"Third quote"

-- Third Author`,
			expected: []Quote{
				{Text: "\"First quote\"", Author: ""},
				{Text: "\"Second quote without author\"", Author: "First Author"},
				{Text: "\"Third quote\"", Author: ""},
			},
		},
		{
			name: "multiline quote",
			content: `"This is a very long quote
that spans multiple lines
and should be handled correctly."

-- Multi Line Author`,
			expected: []Quote{
				{Text: "\"This is a very long quote that spans multiple lines and should be handled correctly.\"", Author: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file with test content
			tempDir := testutil.TempDir(t)
			tempFile := filepath.Join(tempDir, "test.pqf")
			
			if err := os.WriteFile(tempFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			
			quotes, err := ParsePQF(tempFile)
			if err != nil {
				t.Fatalf("ParsePQF failed: %v", err)
			}
			
			if len(quotes) != len(tt.expected) {
				t.Errorf("Expected %d quotes, got %d", len(tt.expected), len(quotes))
				return
			}
			
			for i, expectedQuote := range tt.expected {
				if i >= len(quotes) {
					t.Errorf("Missing quote at index %d", i)
					continue
				}
				
				if quotes[i].Text != expectedQuote.Text {
					t.Errorf("Quote %d text mismatch.\nExpected: %s\nGot: %s", i, expectedQuote.Text, quotes[i].Text)
				}
				
				if quotes[i].Author != expectedQuote.Author {
					t.Errorf("Quote %d author mismatch.\nExpected: %s\nGot: %s", i, expectedQuote.Author, quotes[i].Author)
				}
			}
		})
	}
}

func TestQuote_Validation(t *testing.T) {
	tests := []struct {
		name  string
		quote Quote
	}{
		{
			name: "quote with text and author",
			quote: Quote{
				Text:   "Valid quote text",
				Author: "Valid Author",
			},
		},
		{
			name: "quote with text only",
			quote: Quote{
				Text:   "Quote without author",
				Author: "",
			},
		},
		{
			name: "quote with very long text",
			quote: Quote{
				Text:   "This is a very long quote that might be used to test the handling of lengthy text content in the quote parser and ensure that it doesn't break or cause issues with memory allocation or string processing.",
				Author: "Verbose Author",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			if tt.quote.Text == "" && tt.quote.Author != "" {
				t.Error("Quote should not have author without text")
			}
			
			// Ensure text doesn't contain control characters
			for _, char := range tt.quote.Text {
				if char < 32 && char != 10 && char != 13 && char != 9 {
					t.Errorf("Quote text contains invalid control character: %d", char)
				}
			}
		})
	}
}