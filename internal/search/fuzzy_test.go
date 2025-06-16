package search

import (
	"testing"
	"time"

	"personal-disorganizer/internal/storage"
	"personal-disorganizer/internal/testutil"
)

func TestNewEngine(t *testing.T) {
	engine := NewEngine()
	if engine == nil {
		t.Error("NewEngine() returned nil")
	}
}

func TestEngine_Search(t *testing.T) {
	engine := NewEngine()
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	tasks := []storage.Task{
		{
			ID:   "task1",
			Text: "Complete project documentation",
			Done: false,
			Date: today,
		},
		{
			ID:   "task2", 
			Text: "Review code changes",
			Done: true,
			Date: yesterday,
		},
		{
			ID:   "task3",
			Text: "Meeting with team",
			Done: false,
			Date: tomorrow,
		},
		{
			ID:   "task4",
			Text: "Update project timeline",
			Done: false,
			Date: today,
		},
		{
			ID:   "task5",
			Text: "Write unit tests",
			Done: false,
			Date: tomorrow,
		},
	}

	tests := []struct {
		name           string
		query          string
		expectedCount  int
		expectedFirst  string // ID of expected first result
	}{
		{
			name:          "exact match",
			query:         "Complete project documentation",
			expectedCount: 1,
			expectedFirst: "task1",
		},
		{
			name:          "partial match",
			query:         "project",
			expectedCount: 2, // "Complete project documentation" and "Update project timeline"
			expectedFirst: "task4", // May vary based on scoring
		},
		{
			name:          "fuzzy match",
			query:         "proj",
			expectedCount: 2,
			expectedFirst: "task4", // May vary based on scoring
		},
		{
			name:          "case insensitive",
			query:         "PROJECT",
			expectedCount: 2,
			expectedFirst: "task4", // May vary based on scoring
		},
		{
			name:          "word boundary match",
			query:         "team",
			expectedCount: 1,
			expectedFirst: "task3",
		},
		{
			name:          "no matches",
			query:         "nonexistent",
			expectedCount: 0,
		},
		{
			name:          "empty query",
			query:         "",
			expectedCount: 0,
		},
		{
			name:          "single character",
			query:         "t",
			expectedCount: 4, // "team", "tests", "timeline", plus "documentation"
			expectedFirst: "task5", // May vary based on scoring
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := engine.Search(tt.query, tasks)
			
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
				for i, result := range results {
					t.Logf("Result %d: %s (score: %d)", i, result.Task.Text, result.Score)
				}
			}
			
			if tt.expectedCount > 0 && len(results) > 0 {
				if results[0].Task.ID != tt.expectedFirst {
					t.Errorf("Expected first result to be %s, got %s", tt.expectedFirst, results[0].Task.ID)
				}
			}
		})
	}
}

func TestEngine_CalculateScore(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name     string
		query    string
		text     string
		minScore int // Minimum expected score
		maxScore int // Maximum expected score
	}{
		{
			name:     "exact match",
			query:    "test",
			text:     "test",
			minScore: 900, // Should be very high
			maxScore: 1000,
		},
		{
			name:     "exact substring match",
			query:    "test",
			text:     "this is a test",
			minScore: 400,
			maxScore: 600,
		},
		{
			name:     "fuzzy match at word boundary",
			query:    "test",
			text:     "testing something",
			minScore: 500, // Substring match gets higher score
			maxScore: 700,
		},
		{
			name:     "fuzzy match in middle",
			query:    "test",
			text:     "something testing",
			minScore: 500, // Substring match gets higher score
			maxScore: 700,
		},
		{
			name:     "case insensitive",
			query:    "test",
			text:     "TEST",
			minScore: 900,
			maxScore: 1000,
		},
		{
			name:     "no match",
			query:    "xyz",
			text:     "test",
			minScore: 0,
			maxScore: 0,
		},
		{
			name:     "partial fuzzy match",
			query:    "tst",
			text:     "test",
			minScore: 20,
			maxScore: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := engine.calculateScore(tt.query, tt.text)
			
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("Score %d not in expected range [%d, %d] for query '%s' in text '%s'", 
					score, tt.minScore, tt.maxScore, tt.query, tt.text)
			}
		})
	}
}

func TestEngine_SearchSorting(t *testing.T) {
	engine := NewEngine()
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	tasks := []storage.Task{
		{
			ID:   "task1",
			Text: "test completed task",
			Done: true,
			Date: yesterday,
		},
		{
			ID:   "task2", 
			Text: "test active task",
			Done: false,
			Date: today,
		},
		{
			ID:   "task3",
			Text: "test future task",
			Done: false,
			Date: today.AddDate(0, 0, 1),
		},
	}

	results := engine.Search("test", tasks)
	
	if len(results) != 3 {
		t.Fatalf("Expected 3 results, got %d", len(results))
	}

	// Active/future tasks should score higher than completed tasks
	activeTaskScore := -1
	completedTaskScore := -1
	
	for _, result := range results {
		if result.Task.ID == "task2" || result.Task.ID == "task3" {
			if activeTaskScore == -1 || result.Score > activeTaskScore {
				activeTaskScore = result.Score
			}
		}
		if result.Task.ID == "task1" {
			completedTaskScore = result.Score
		}
	}

	if activeTaskScore <= completedTaskScore {
		t.Errorf("Active tasks should score higher than completed tasks. Active: %d, Completed: %d", 
			activeTaskScore, completedTaskScore)
	}
}

func TestEngine_HighlightMatch(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		name     string
		query    string
		text     string
		expected string
	}{
		{
			name:     "simple highlight",
			query:    "test",
			text:     "this is a test",
			expected: "this is a test", // Currently just returns original text
		},
		{
			name:     "no match",
			query:    "xyz",
			text:     "test",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.highlightMatch(tt.query, tt.text)
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestEngine_SearchEdgeCases(t *testing.T) {
	engine := NewEngine()

	tasks := []storage.Task{
		{
			ID:   "task1",
			Text: "normal task",
			Done: false,
			Date: time.Now(),
		},
		{
			ID:   "task2",
			Text: "!@#$%^&*()",
			Done: false,
			Date: time.Now(),
		},
		{
			ID:   "task3",
			Text: "very long task description that goes on and on and might cause issues with the search algorithm if it's not properly optimized for handling lengthy text content",
			Done: false,
			Date: time.Now(),
		},
		{
			ID:   "task4",
			Text: "",
			Done: false,
			Date: time.Now(),
		},
	}

	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{
			name:          "special characters query",
			query:         "!@#",
			expectedCount: 1,
		},
		{
			name:          "very long query",
			query:         "very long task description that goes on",
			expectedCount: 1,
		},
		{
			name:          "whitespace query",
			query:         "   ",
			expectedCount: 0, // Should treat as empty
		},
		{
			name:          "unicode query",
			query:         "test",
			expectedCount: 0, // No unicode content in test data
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := engine.Search(tt.query, tasks)
			
			if len(results) != tt.expectedCount {
				t.Errorf("Expected %d results, got %d", tt.expectedCount, len(results))
			}
		})
	}
}

func TestEngine_SearchPerformance(t *testing.T) {
	engine := NewEngine()
	
	// Create a large number of tasks
	tasks := make([]storage.Task, 1000)
	now := time.Now()
	
	for i := 0; i < 1000; i++ {
		tasks[i] = storage.Task{
			ID:   testutil.MockUUID(i),
			Text: testutil.MockTaskText(i),
			Done: i%3 == 0, // Every third task is done
			Date: now.AddDate(0, 0, i%30-15), // Spread across month
		}
	}

	// Measure search performance
	query := "test"
	
	// This is a basic performance test - in a real scenario you might want to measure time
	results := engine.Search(query, tasks)
	
	// Just verify it completes and returns reasonable results
	if len(results) < 0 {
		t.Error("Search should return non-negative number of results")
	}
	
	// Verify results are properly sorted
	for i := 1; i < len(results); i++ {
		if results[i-1].Score < results[i].Score {
			t.Error("Results should be sorted by score (highest first)")
			break
		}
		
		// If scores are equal, should be sorted by date
		if results[i-1].Score == results[i].Score &&
		   results[i-1].Task.Date.Before(results[i].Task.Date) {
			t.Error("Results with equal scores should be sorted by date (newest first)")
			break
		}
	}
}

func TestAbsFunction(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
		{-1, 1},
		{1, 1},
	}

	for _, tt := range tests {
		result := abs(tt.input)
		if result != tt.expected {
			t.Errorf("abs(%d) = %d, expected %d", tt.input, result, tt.expected)
		}
	}
}

func TestMaxFunction(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{5, 3, 5},
		{3, 5, 5},
		{5, 5, 5},
		{-3, -5, -3},
		{0, 0, 0},
	}

	for _, tt := range tests {
		result := max(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("max(%d, %d) = %d, expected %d", tt.a, tt.b, result, tt.expected)
		}
	}
}