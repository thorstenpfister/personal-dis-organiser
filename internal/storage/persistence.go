package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// Config represents the application configuration
type Config struct {
	CalendarURLs    []string `json:"calendar_urls"`
	DataFile        string   `json:"data_file"`
	QuoteFiles      []string `json:"quote_files"`
	RefreshInterval int      `json:"refresh_interval"`
	DateFormat      string   `json:"date_format"`
	TimeFormat      string   `json:"time_format"`
	Theme           string   `json:"theme"`
}

// Task represents a single task or calendar event
type Task struct {
	ID         string    `json:"id"`
	Text       string    `json:"text"`
	Done       bool      `json:"done"`
	Date       time.Time `json:"date"`
	IsCalendar bool      `json:"is_calendar"`
	StartTime  time.Time `json:"start_time"`
	Priority   int       `json:"priority"`
	CreatedAt  time.Time `json:"created_at"`
	Level      int       `json:"level"` // Hierarchy level (0 = top level)
}

// AppData represents all application data
type AppData struct {
	Tasks    []Task    `json:"tasks"`
	Settings Settings  `json:"settings"`
}

// Settings represents application settings
type Settings struct {
	LastQuoteIndex       int `json:"last_quote_index"`
	TasksCompletedToday  int `json:"tasks_completed_today"`
}

// Storage handles data persistence
type Storage struct {
	configDir string
	dataPath  string
	config    *Config
}

// NewStorage creates a new storage instance
func NewStorage() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	
	configDir := filepath.Join(homeDir, ".config", "personal-disorganizer")
	dataPath := filepath.Join(configDir, "data.json")
	
	s := &Storage{
		configDir: configDir,
		dataPath:  dataPath,
	}
	
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}
	
	// Load or create config
	if err := s.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	
	return s, nil
}

// loadConfig loads configuration or creates default config
func (s *Storage) loadConfig() error {
	configPath := filepath.Join(s.configDir, "config.json")
	
	// Create default config if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := &Config{
			CalendarURLs:    []string{},
			DataFile:        "data.json",
			QuoteFiles:      []string{},
			RefreshInterval: 300,
			DateFormat:      "2006-01-02",
			TimeFormat:      "15:04",
			Theme:           "dracula",
		}
		
		if err := s.saveConfig(defaultConfig); err != nil {
			return fmt.Errorf("failed to save default config: %w", err)
		}
		
		s.config = defaultConfig
		return nil
	}
	
	// Load existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}
	
	s.config = config
	return nil
}

// saveConfig saves configuration to file
func (s *Storage) saveConfig(config *Config) error {
	configPath := filepath.Join(s.configDir, "config.json")
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// LoadData loads application data from file
func (s *Storage) LoadData() (*AppData, error) {
	// Create default data if file doesn't exist
	if _, err := os.Stat(s.dataPath); os.IsNotExist(err) {
		defaultData := &AppData{
			Tasks: []Task{},
			Settings: Settings{
				LastQuoteIndex:      0,
				TasksCompletedToday: 0,
			},
		}
		return defaultData, nil
	}
	
	data, err := os.ReadFile(s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}
	
	appData := &AppData{}
	if err := json.Unmarshal(data, appData); err != nil {
		return nil, fmt.Errorf("failed to parse data file: %w", err)
	}
	
	return appData, nil
}

// SaveData saves application data to file
func (s *Storage) SaveData(data *AppData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}
	
	if err := os.WriteFile(s.dataPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}
	
	return nil
}

// GetConfig returns the current configuration
func (s *Storage) GetConfig() *Config {
	return s.config
}

// CreateTask creates a new task with a unique ID
func (s *Storage) CreateTask(text string, date time.Time) *Task {
	return &Task{
		ID:        uuid.New().String(),
		Text:      text,
		Done:      false,
		Date:      date,
		IsCalendar: false,
		Priority:  0,
		CreatedAt: time.Now(),
		Level:     0,
	}
}

// LogError logs an error to the error log file
func (s *Storage) LogError(err error) {
	logPath := filepath.Join(s.configDir, "error.log")
	
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] ERROR: %s\n", timestamp, err.Error())
	
	file, openErr := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return // Can't log the error if we can't open the log file
	}
	defer file.Close()
	
	file.WriteString(logEntry)
}

// PurgeData deletes all application data and config files
func (s *Storage) PurgeData() error {
	// Remove the entire config directory and all its contents
	if err := os.RemoveAll(s.configDir); err != nil {
		return fmt.Errorf("failed to remove config directory: %w", err)
	}
	
	return nil
}