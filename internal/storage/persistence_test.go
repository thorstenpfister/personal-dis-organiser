package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"personal-disorganizer/internal/testutil"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful creation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use temporary directory for testing
			originalHome := os.Getenv("HOME")
			tempDir := testutil.TempDir(t)
			os.Setenv("HOME", tempDir)
			defer os.Setenv("HOME", originalHome)

			storage, err := NewStorage()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if storage == nil {
					t.Error("NewStorage() returned nil storage")
				}

				// Verify config directory was created
				expectedConfigDir := filepath.Join(tempDir, ".config", "personal-disorganizer")
				testutil.AssertFileExists(t, expectedConfigDir)

				// Verify default config was created
				configPath := filepath.Join(expectedConfigDir, "config.json")
				testutil.AssertFileExists(t, configPath)

				// Verify config content
				config := storage.GetConfig()
				if config == nil {
					t.Error("Config is nil")
				}
				if config.Theme != "dracula" {
					t.Errorf("Expected default theme 'dracula', got %s", config.Theme)
				}
			}
		})
	}
}

func TestStorage_LoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		setupConfig    func(dir string) error
		expectedTheme  string
		expectedURLs   int
		expectError    bool
	}{
		{
			name: "load existing valid config",
			setupConfig: func(dir string) error {
				configData := &Config{
					CalendarURLs:    []string{"https://example.com/calendar.ics"},
					DataFile:        "data.json",
					QuoteFiles:      []string{"quotes/test.json"},
					RefreshInterval: 300,
					DateFormat:      "2006-01-02",
					TimeFormat:      "15:04",
					Theme:           "light",
				}
				testutil.CreateTestConfig(t, dir, configData)
				return nil
			},
			expectedTheme: "light",
			expectedURLs:  1,
			expectError:   false,
		},
		{
			name: "create default config when none exists",
			setupConfig: func(dir string) error {
				// Don't create any config file
				return nil
			},
			expectedTheme: "dracula",
			expectedURLs:  0,
			expectError:   false,
		},
		{
			name: "handle corrupted config file",
			setupConfig: func(dir string) error {
				configPath := filepath.Join(dir, "config.json")
				return os.WriteFile(configPath, []byte("invalid json"), 0644)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := testutil.TempDir(t)
			
			// Setup test config
			if err := tt.setupConfig(tempDir); err != nil {
				t.Fatalf("Failed to setup test config: %v", err)
			}

			storage := &Storage{
				configDir: tempDir,
				dataPath:  filepath.Join(tempDir, "data.json"),
			}

			err := storage.loadConfig()
			
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

			config := storage.GetConfig()
			if config.Theme != tt.expectedTheme {
				t.Errorf("Expected theme %s, got %s", tt.expectedTheme, config.Theme)
			}
			if len(config.CalendarURLs) != tt.expectedURLs {
				t.Errorf("Expected %d calendar URLs, got %d", tt.expectedURLs, len(config.CalendarURLs))
			}
		})
	}
}

func TestStorage_LoadData(t *testing.T) {
	tests := []struct {
		name         string
		setupData    func(dir string) error
		expectedTasks int
		expectError  bool
	}{
		{
			name: "load existing valid data",
			setupData: func(dir string) error {
				now := time.Now()
				data := &AppData{
					Tasks: []Task{
						{
							ID:         "test-task-1",
							Text:       "Test task 1",
							Done:       false,
							Date:       now,
							IsCalendar: false,
							Priority:   0,
							CreatedAt:  now,
							Level:      0,
						},
						{
							ID:         "test-task-2",
							Text:       "Test task 2",
							Done:       true,
							Date:       now.AddDate(0, 0, -1),
							IsCalendar: false,
							Priority:   1,
							CreatedAt:  now.AddDate(0, 0, -1),
							Level:      0,
						},
					},
					Settings: Settings{
						LastQuoteIndex:      5,
						TasksCompletedToday: 2,
					},
				}
				testutil.CreateTestData(t, dir, data)
				return nil
			},
			expectedTasks: 2,
			expectError:   false,
		},
		{
			name: "create default data when none exists",
			setupData: func(dir string) error {
				// Don't create any data file
				return nil
			},
			expectedTasks: 0,
			expectError:   false,
		},
		{
			name: "handle corrupted data file",
			setupData: func(dir string) error {
				dataPath := filepath.Join(dir, "data.json")
				return os.WriteFile(dataPath, []byte("invalid json"), 0644)
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := testutil.TempDir(t)
			
			// Setup test data
			if err := tt.setupData(tempDir); err != nil {
				t.Fatalf("Failed to setup test data: %v", err)
			}

			storage := &Storage{
				configDir: tempDir,
				dataPath:  filepath.Join(tempDir, "data.json"),
			}

			data, err := storage.LoadData()
			
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

			if len(data.Tasks) != tt.expectedTasks {
				t.Errorf("Expected %d tasks, got %d", tt.expectedTasks, len(data.Tasks))
			}
		})
	}
}

func TestStorage_SaveData(t *testing.T) {
	tempDir := testutil.TempDir(t)
	dataPath := filepath.Join(tempDir, "data.json")

	storage := &Storage{
		configDir: tempDir,
		dataPath:  dataPath,
	}

	now := time.Now()
	testData := &AppData{
		Tasks: []Task{
			{
				ID:         "test-task-1",
				Text:       "Test task 1",
				Done:       false,
				Date:       now,
				IsCalendar: false,
				Priority:   0,
				CreatedAt:  now,
				Level:      0,
			},
		},
		Settings: Settings{
			LastQuoteIndex:      5,
			TasksCompletedToday: 1,
		},
	}
	
	// Test saving data
	err := storage.SaveData(testData)
	if err != nil {
		t.Errorf("SaveData() error = %v", err)
		return
	}

	// Verify file was created
	testutil.AssertFileExists(t, dataPath)

	// Verify file content
	fileData, err := os.ReadFile(dataPath)
	if err != nil {
		t.Fatalf("Failed to read saved data file: %v", err)
	}

	var savedData AppData
	if err := json.Unmarshal(fileData, &savedData); err != nil {
		t.Fatalf("Failed to parse saved data: %v", err)
	}

	if len(savedData.Tasks) != len(testData.Tasks) {
		t.Errorf("Expected %d tasks, got %d", len(testData.Tasks), len(savedData.Tasks))
	}

	if savedData.Settings.TasksCompletedToday != testData.Settings.TasksCompletedToday {
		t.Errorf("Expected tasks completed today %d, got %d", 
			testData.Settings.TasksCompletedToday, savedData.Settings.TasksCompletedToday)
	}
}

func TestStorage_CreateTask(t *testing.T) {
	storage := &Storage{}
	
	taskText := "Test task"
	taskDate := testutil.FixedTime()
	
	task := storage.CreateTask(taskText, taskDate)
	
	if task == nil {
		t.Error("CreateTask() returned nil")
		return
	}
	
	if task.Text != taskText {
		t.Errorf("Expected task text %s, got %s", taskText, task.Text)
	}
	
	if !task.Date.Equal(taskDate) {
		t.Errorf("Expected task date %v, got %v", taskDate, task.Date)
	}
	
	if task.Done {
		t.Error("New task should not be done")
	}
	
	if task.IsCalendar {
		t.Error("New task should not be calendar task")
	}
	
	if task.ID == "" {
		t.Error("Task ID should not be empty")
	}
	
	if task.Level != 0 {
		t.Errorf("Expected task level 0, got %d", task.Level)
	}
}

func TestStorage_LogError(t *testing.T) {
	tempDir := testutil.TempDir(t)
	
	storage := &Storage{
		configDir: tempDir,
	}
	
	testError := testutil.MockError("test error message")
	
	// Test logging error
	storage.LogError(testError)
	
	// Verify log file was created
	logPath := filepath.Join(tempDir, "error.log")
	testutil.AssertFileExists(t, logPath)
	
	// Verify log content
	logData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	
	logContent := string(logData)
	if !strings.Contains(logContent, "test error message") {
		t.Errorf("Log content should contain error message, got: %s", logContent)
	}
	
	if !strings.Contains(logContent, "ERROR:") {
		t.Errorf("Log content should contain ERROR prefix, got: %s", logContent)
	}
}

func TestStorage_PurgeData(t *testing.T) {
	tempDir := testutil.TempDir(t)
	
	// Create some test files in the config directory
	configFile := filepath.Join(tempDir, "config.json")
	dataFile := filepath.Join(tempDir, "data.json")
	logFile := filepath.Join(tempDir, "error.log")
	
	if err := os.WriteFile(configFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	if err := os.WriteFile(dataFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create test data file: %v", err)
	}
	if err := os.WriteFile(logFile, []byte("test log"), 0644); err != nil {
		t.Fatalf("Failed to create test log file: %v", err)
	}
	
	// Verify files exist before purge
	testutil.AssertFileExists(t, configFile)
	testutil.AssertFileExists(t, dataFile)
	testutil.AssertFileExists(t, logFile)
	
	storage := &Storage{
		configDir: tempDir,
	}
	
	// Test purge operation
	err := storage.PurgeData()
	if err != nil {
		t.Errorf("PurgeData() error = %v", err)
		return
	}
	
	// Verify all files were removed
	testutil.AssertFileNotExists(t, configFile)
	testutil.AssertFileNotExists(t, dataFile)
	testutil.AssertFileNotExists(t, logFile)
	testutil.AssertFileNotExists(t, tempDir)
}

func TestTask_Validation(t *testing.T) {
	tests := []struct {
		name string
		task Task
		want bool
	}{
		{
			name: "valid regular task",
			task: Task{
				ID:         "test-1",
				Text:       "Test task",
				Done:       false,
				Date:       time.Now(),
				IsCalendar: false,
				Priority:   0,
				CreatedAt:  time.Now(),
				Level:      0,
			},
			want: true,
		},
		{
			name: "valid calendar task",
			task: Task{
				ID:         "cal-1",
				Text:       "Meeting",
				Done:       false,
				Date:       time.Now(),
				IsCalendar: true,
				StartTime:  time.Now(),
				Priority:   -1,
				CreatedAt:  time.Now(),
				Level:      0,
			},
			want: true,
		},
		{
			name: "task with hierarchical level",
			task: Task{
				ID:         "test-2",
				Text:       "Subtask",
				Done:       false,
				Date:       time.Now(),
				IsCalendar: false,
				Priority:   0,
				CreatedAt:  time.Now(),
				Level:      1,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			if tt.task.ID == "" {
				t.Error("Task ID should not be empty")
			}
			if tt.task.Text == "" {
				t.Error("Task text should not be empty")
			}
			if tt.task.Level < 0 {
				t.Error("Task level should not be negative")
			}
		})
	}
}