package calendar

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"personal-disorganizer/internal/storage"
)

// Event represents a calendar event
type Event struct {
	Summary     string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Location    string
}

// Logger interface for error logging
type Logger interface {
	LogError(err error)
}

// Manager handles calendar integration
type Manager struct {
	urls   []string
	logger Logger
}

// NewManager creates a new calendar manager
func NewManager(urls []string) *Manager {
	return &Manager{
		urls: urls,
	}
}

// SetLogger sets the logger instance for error logging
func (m *Manager) SetLogger(logger Logger) {
	m.logger = logger
}

// FetchEvents fetches events from all configured calendars for a specific date
func (m *Manager) FetchEvents(date time.Time) ([]storage.Task, error) {
	var allTasks []storage.Task
	
	for _, url := range m.urls {
		events, err := m.fetchEventsFromURL(url, date)
		if err != nil {
			// Log error but continue with other calendars
			if m.logger != nil {
				m.logger.LogError(fmt.Errorf("calendar fetch failed for %s: %w", url, err))
			}
			continue
		}
		
		// Convert events to tasks
		for _, event := range events {
			task := storage.Task{
				ID:         fmt.Sprintf("cal_%d", time.Now().UnixNano()),
				Text:       event.Summary,
				Done:       false,
				Date:       date,
				IsCalendar: true,
				StartTime:  event.StartTime,
				Priority:   -1, // Calendar events have highest priority
				CreatedAt:  time.Now(),
				Level:      0,
			}
			allTasks = append(allTasks, task)
		}
	}
	
	return allTasks, nil
}

// fetchEventsFromURL fetches events from a single iCal URL
func (m *Manager) fetchEventsFromURL(url string, date time.Time) ([]Event, error) {
	// Handle webcal:// URLs
	if strings.HasPrefix(url, "webcal://") {
		url = "https://" + url[9:]
	}
	
	// Fetch the iCal data
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch calendar: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("calendar request failed: %d", resp.StatusCode)
	}
	
	// Parse the iCal data
	events, err := m.parseICalData(resp.Body, date)
	if err != nil && m.logger != nil {
		m.logger.LogError(fmt.Errorf("calendar parse failed for %s: %w", url, err))
	}
	return events, err
}

// parseICalData parses iCal data and extracts events for the specified date
func (m *Manager) parseICalData(reader io.Reader, targetDate time.Time) ([]Event, error) {
	var events []Event
	var currentEvent *Event
	
	scanner := bufio.NewScanner(reader)
	
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		if line == "BEGIN:VEVENT" {
			currentEvent = &Event{}
		} else if line == "END:VEVENT" {
			if currentEvent != nil {
				// Check if event occurs on target date
				if m.eventOccursOnDate(*currentEvent, targetDate) {
					events = append(events, *currentEvent)
				}
			}
			currentEvent = nil
		} else if currentEvent != nil {
			m.parseEventLine(currentEvent, line)
		}
	}
	
	return events, scanner.Err()
}

// parseEventLine parses a single line of event data
func (m *Manager) parseEventLine(event *Event, line string) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return
	}
	
	key := strings.ToUpper(parts[0])
	value := parts[1]
	
	switch {
	case strings.HasPrefix(key, "SUMMARY"):
		event.Summary = value
	case strings.HasPrefix(key, "DESCRIPTION"):
		event.Description = value
	case strings.HasPrefix(key, "LOCATION"):
		event.Location = value
	case strings.HasPrefix(key, "DTSTART"):
		if t, err := m.parseDateTime(value); err == nil {
			event.StartTime = t
		}
	case strings.HasPrefix(key, "DTEND"):
		if t, err := m.parseDateTime(value); err == nil {
			event.EndTime = t
		}
	}
}

// parseDateTime parses iCal datetime format
func (m *Manager) parseDateTime(value string) (time.Time, error) {
	// Remove timezone info for now - simplified parsing
	value = strings.Split(value, ";")[0]
	
	// Try different datetime formats
	formats := []string{
		"20060102T150405Z",
		"20060102T150405",
		"20060102",
	}
	
	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			return t, nil
		}
	}
	
	return time.Time{}, fmt.Errorf("unable to parse datetime: %s", value)
}

// eventOccursOnDate checks if an event occurs on the specified date
func (m *Manager) eventOccursOnDate(event Event, date time.Time) bool {
	eventDate := event.StartTime.Truncate(24 * time.Hour)
	targetDate := date.Truncate(24 * time.Hour)
	return eventDate.Equal(targetDate)
}