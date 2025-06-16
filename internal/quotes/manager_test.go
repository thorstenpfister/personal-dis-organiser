package quotes

import (
	"os"
	"path/filepath"
	"testing"

	"personal-disorganizer/internal/parser"
	"personal-disorganizer/internal/testutil"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name       string
		configDir  string
		quoteFiles []string
		expectErr  bool
	}{
		{
			name:       "create manager with no quote files",
			configDir:  "/tmp/test",
			quoteFiles: []string{},
			expectErr:  false,
		},
		{
			name:       "create manager with valid quote files",
			configDir:  "testdata",
			quoteFiles: []string{"test1.json"},
			expectErr:  false,
		},
		{
			name:       "create manager with missing quote files",
			configDir:  "testdata",
			quoteFiles: []string{"nonexistent.json"},
			expectErr:  false, // Should not error, just skip missing files
		},
		{
			name:       "create manager with mixed files",
			configDir:  "testdata",
			quoteFiles: []string{"test1.json", "nonexistent.json", "test2.json"},
			expectErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.configDir, tt.quoteFiles)
			
			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if manager == nil {
				t.Error("Manager should not be nil")
			}
		})
	}
}

func TestManager_LoadQuoteFile(t *testing.T) {
	tempDir := testutil.TempDir(t)
	
	// Create test quote file
	testQuotes := `[
		{
			"text": "Test quote",
			"author": "Test Author"
		}
	]`
	testFile := filepath.Join(tempDir, "test.json")
	if err := os.WriteFile(testFile, []byte(testQuotes), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		configDir   string
		filename    string
		expectErr   bool
		expectCount int
	}{
		{
			name:        "load relative path",
			configDir:   tempDir,
			filename:    "test.json",
			expectErr:   false,
			expectCount: 1,
		},
		{
			name:        "load absolute path",
			configDir:   "/tmp",
			filename:    testFile,
			expectErr:   false,
			expectCount: 1,
		},
		{
			name:        "load nonexistent file",
			configDir:   tempDir,
			filename:    "nonexistent.json",
			expectErr:   true,
			expectCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				configDir: tt.configDir,
				quotes:    []parser.Quote{},
			}
			
			err := manager.loadQuoteFile(tt.filename)
			
			if tt.expectErr && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if len(manager.quotes) != tt.expectCount {
				t.Errorf("Expected %d quotes, got %d", tt.expectCount, len(manager.quotes))
			}
		})
	}
}

func TestManager_GetRandomQuote(t *testing.T) {
	tests := []struct {
		name         string
		setupQuotes  func() *Manager
		expectNil    bool
		expectQuote  bool
	}{
		{
			name: "get quote from populated manager",
			setupQuotes: func() *Manager {
				return &Manager{
					quotes: []parser.Quote{
						{Text: "Quote 1", Author: "Author 1"},
						{Text: "Quote 2", Author: "Author 2"},
						{Text: "Quote 3", Author: "Author 3"},
					},
				}
			},
			expectNil:   false,
			expectQuote: true,
		},
		{
			name: "get quote from empty manager",
			setupQuotes: func() *Manager {
				return &Manager{
					quotes: []parser.Quote{},
				}
			},
			expectNil:   true,
			expectQuote: false,
		},
		{
			name: "get quote from nil quotes",
			setupQuotes: func() *Manager {
				return &Manager{
					quotes: nil,
				}
			},
			expectNil:   true,
			expectQuote: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := tt.setupQuotes()
			
			quote := manager.GetRandomQuote()
			
			if tt.expectNil && quote != nil {
				t.Error("Expected nil quote but got one")
			}
			
			if tt.expectQuote && quote == nil {
				t.Error("Expected quote but got nil")
			}
			
			if quote != nil {
				if quote.Text == "" {
					t.Error("Quote text should not be empty")
				}
			}
		})
	}
}

func TestManager_GetRandomQuote_Randomness(t *testing.T) {
	manager := &Manager{
		quotes: []parser.Quote{
			{Text: "Quote 1", Author: "Author 1"},
			{Text: "Quote 2", Author: "Author 2"},
			{Text: "Quote 3", Author: "Author 3"},
			{Text: "Quote 4", Author: "Author 4"},
			{Text: "Quote 5", Author: "Author 5"},
		},
	}

	// Get multiple quotes and check for some variation
	quotes := make(map[string]bool)
	iterations := 20
	
	for i := 0; i < iterations; i++ {
		quote := manager.GetRandomQuote()
		if quote != nil {
			quotes[quote.Text] = true
		}
	}

	// With 5 quotes and 20 iterations, we should see some variety
	// This is probabilistic, but should pass most of the time
	if len(quotes) < 2 {
		t.Errorf("Expected some randomness in quote selection, got %d unique quotes out of %d iterations", 
			len(quotes), iterations)
	}
}

func TestManager_GetQuoteCount(t *testing.T) {
	tests := []struct {
		name          string
		quotes        []parser.Quote
		expectedCount int
	}{
		{
			name:          "empty quotes",
			quotes:        []parser.Quote{},
			expectedCount: 0,
		},
		{
			name: "single quote",
			quotes: []parser.Quote{
				{Text: "Single quote", Author: "Single author"},
			},
			expectedCount: 1,
		},
		{
			name: "multiple quotes",
			quotes: []parser.Quote{
				{Text: "Quote 1", Author: "Author 1"},
				{Text: "Quote 2", Author: "Author 2"},
				{Text: "Quote 3", Author: "Author 3"},
			},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				quotes: tt.quotes,
			}
			
			count := manager.GetQuoteCount()
			
			if count != tt.expectedCount {
				t.Errorf("Expected count %d, got %d", tt.expectedCount, count)
			}
		})
	}
}

func TestManager_HasQuotes(t *testing.T) {
	tests := []struct {
		name     string
		quotes   []parser.Quote
		expected bool
	}{
		{
			name:     "empty quotes",
			quotes:   []parser.Quote{},
			expected: false,
		},
		{
			name: "with quotes",
			quotes: []parser.Quote{
				{Text: "Quote", Author: "Author"},
			},
			expected: true,
		},
		{
			name:     "nil quotes",
			quotes:   nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{
				quotes: tt.quotes,
			}
			
			hasQuotes := manager.HasQuotes()
			
			if hasQuotes != tt.expected {
				t.Errorf("Expected HasQuotes to return %v, got %v", tt.expected, hasQuotes)
			}
		})
	}
}

func TestManager_IntegrationTest(t *testing.T) {
	// Test loading quotes from actual files
	configDir := "testdata"
	quoteFiles := []string{"test1.json", "test2.json"}
	
	manager, err := NewManager(configDir, quoteFiles)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	
	// Should have loaded quotes from both files
	expectedCount := 5 // 2 from test1.json + 3 from test2.json
	if manager.GetQuoteCount() != expectedCount {
		t.Errorf("Expected %d quotes, got %d", expectedCount, manager.GetQuoteCount())
	}
	
	// Should have quotes available
	if !manager.HasQuotes() {
		t.Error("Manager should have quotes")
	}
	
	// Should be able to get a random quote
	quote := manager.GetRandomQuote()
	if quote == nil {
		t.Error("Should be able to get a random quote")
	}
	
	if quote != nil && quote.Text == "" {
		t.Error("Quote text should not be empty")
	}
}

func TestManager_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func() (*Manager, error)
		shouldWork bool
	}{
		{
			name: "invalid JSON file",
			setupFunc: func() (*Manager, error) {
				return NewManager("testdata", []string{"invalid.json"})
			},
			shouldWork: true, // Should create manager but skip invalid file
		},
		{
			name: "empty JSON file",
			setupFunc: func() (*Manager, error) {
				return NewManager("testdata", []string{"empty.json"})
			},
			shouldWork: true,
		},
		{
			name: "mixed valid and invalid files",
			setupFunc: func() (*Manager, error) {
				return NewManager("testdata", []string{"test1.json", "invalid.json", "test2.json"})
			},
			shouldWork: true, // Should load valid files and skip invalid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := tt.setupFunc()
			
			if tt.shouldWork {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if manager == nil {
					t.Error("Manager should not be nil")
				}
			} else {
				if err == nil {
					t.Error("Expected error but got none")
				}
			}
		})
	}
}

func TestManager_PathHandling(t *testing.T) {
	tempDir := testutil.TempDir(t)
	
	// Create nested directory structure
	quotesDir := filepath.Join(tempDir, "quotes")
	if err := os.MkdirAll(quotesDir, 0755); err != nil {
		t.Fatalf("Failed to create quotes directory: %v", err)
	}
	
	// Create test quote file in nested directory
	testQuotes := `[{"text": "Nested quote", "author": "Nested author"}]`
	nestedFile := filepath.Join(quotesDir, "nested.json")
	if err := os.WriteFile(nestedFile, []byte(testQuotes), 0644); err != nil {
		t.Fatalf("Failed to create nested test file: %v", err)
	}

	tests := []struct {
		name        string
		configDir   string
		filename    string
		expectCount int
	}{
		{
			name:        "relative path within config dir",
			configDir:   tempDir,
			filename:    "quotes/nested.json",
			expectCount: 1,
		},
		{
			name:        "absolute path",
			configDir:   "/tmp",
			filename:    nestedFile,
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.configDir, []string{tt.filename})
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if manager.GetQuoteCount() != tt.expectCount {
				t.Errorf("Expected %d quotes, got %d", tt.expectCount, manager.GetQuoteCount())
			}
		})
	}
}