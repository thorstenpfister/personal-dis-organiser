package testutil

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TempDir creates a temporary directory for testing
func TempDir(t *testing.T) string {
	t.Helper()
	
	dir, err := os.MkdirTemp("", "personal-disorganizer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})
	
	return dir
}

// CreateTestConfig creates a test configuration file with any interface{}
func CreateTestConfig(t *testing.T, dir string, config interface{}) string {
	t.Helper()
	
	configPath := filepath.Join(dir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	
	return configPath
}

// CreateTestData creates a test data file with any interface{}
func CreateTestData(t *testing.T, dir string, data interface{}) string {
	t.Helper()
	
	dataPath := filepath.Join(dir, "data.json")
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal data: %v", err)
	}
	
	if err := os.WriteFile(dataPath, jsonData, 0644); err != nil {
		t.Fatalf("Failed to write data file: %v", err)
	}
	
	return dataPath
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", path)
	}
}

// AssertFileNotExists checks if a file does not exist
func AssertFileNotExists(t *testing.T, path string) {
	t.Helper()
	
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("Expected file to not exist: %s", path)
	}
}

// AssertJSONEqual compares two JSON-serializable objects
func AssertJSONEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	
	expectedJSON, err := json.MarshalIndent(expected, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal expected: %v", err)
	}
	
	actualJSON, err := json.MarshalIndent(actual, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal actual: %v", err)
	}
	
	if string(expectedJSON) != string(actualJSON) {
		t.Errorf("JSON objects not equal:\nExpected:\n%s\nActual:\n%s", 
			expectedJSON, actualJSON)
	}
}

// FixedTime returns a fixed time for consistent testing
func FixedTime() time.Time {
	return time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
}

// FixedTimeString returns a fixed time as string
func FixedTimeString() string {
	return FixedTime().Format("2006-01-02T15:04:05Z")
}

// MockError creates a mock error for testing
func MockError(message string) error {
	return &mockError{message: message}
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

// MockUUID generates a mock UUID for testing
func MockUUID(id int) string {
	return fmt.Sprintf("test-uuid-%d", id)
}

// MockTaskText generates mock task text for testing
func MockTaskText(id int) string {
	texts := []string{
		"Complete project documentation",
		"Review code changes", 
		"Meeting with team",
		"Update project timeline",
		"Write unit tests",
		"Fix bug in authentication",
		"Deploy to staging environment",
		"Create user manual",
		"Performance optimization",
		"Database migration",
	}
	return texts[id%len(texts)]
}