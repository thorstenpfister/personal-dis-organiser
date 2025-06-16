package testutil

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// MockLogger implements the Logger interface for testing
type MockLogger struct {
	LoggedErrors []error
}

// LogError records errors for testing verification
func (m *MockLogger) LogError(err error) {
	m.LoggedErrors = append(m.LoggedErrors, err)
}

// GetLastError returns the most recent logged error
func (m *MockLogger) GetLastError() error {
	if len(m.LoggedErrors) == 0 {
		return nil
	}
	return m.LoggedErrors[len(m.LoggedErrors)-1]
}

// GetErrorCount returns the number of logged errors
func (m *MockLogger) GetErrorCount() int {
	return len(m.LoggedErrors)
}

// Clear clears all logged errors
func (m *MockLogger) Clear() {
	m.LoggedErrors = nil
}

// MockHTTPClient provides a mock HTTP client for testing
type MockHTTPClient struct {
	Responses map[string]*http.Response
	Errors    map[string]error
}

// NewMockHTTPClient creates a new mock HTTP client
func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{
		Responses: make(map[string]*http.Response),
		Errors:    make(map[string]error),
	}
}

// SetResponse sets a mock response for a URL
func (m *MockHTTPClient) SetResponse(url string, statusCode int, body string) {
	m.Responses[url] = &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

// SetError sets a mock error for a URL
func (m *MockHTTPClient) SetError(url string, err error) {
	m.Errors[url] = err
}

// Get simulates an HTTP GET request
func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	if err, exists := m.Errors[url]; exists {
		return nil, err
	}
	
	if resp, exists := m.Responses[url]; exists {
		return resp, nil
	}
	
	// Default response for unknown URLs
	return &http.Response{
		StatusCode: 404,
		Body:       io.NopCloser(strings.NewReader("Not Found")),
		Header:     make(http.Header),
	}, nil
}

// TimeProvider interface for mockable time
type TimeProvider interface {
	Now() time.Time
}

// MockTimeProvider provides controllable time for testing
type MockTimeProvider struct {
	CurrentTime time.Time
}

// NewMockTimeProvider creates a new mock time provider
func NewMockTimeProvider(t time.Time) *MockTimeProvider {
	return &MockTimeProvider{CurrentTime: t}
}

// Now returns the mock current time
func (m *MockTimeProvider) Now() time.Time {
	return m.CurrentTime
}

// SetTime sets the mock current time
func (m *MockTimeProvider) SetTime(t time.Time) {
	m.CurrentTime = t
}

// Advance advances the mock time by the given duration
func (m *MockTimeProvider) Advance(d time.Duration) {
	m.CurrentTime = m.CurrentTime.Add(d)
}

// Sample iCal data for testing
const SampleICalData = `BEGIN:VCALENDAR
VERSION:2.0
PRODID:Test Calendar
BEGIN:VEVENT
UID:test-event-1@example.com
DTSTART:20240115T100000Z
DTEND:20240115T110000Z
SUMMARY:Test Meeting
DESCRIPTION:This is a test meeting
LOCATION:Conference Room A
END:VEVENT
BEGIN:VEVENT
UID:test-event-2@example.com
DTSTART:20240116T140000Z
DTEND:20240116T150000Z
SUMMARY:Another Meeting
DESCRIPTION:Another test meeting
END:VEVENT
END:VCALENDAR`

// Sample Terry Pratchett Quote File (PQF format) for testing
const SamplePQFData = `"The trouble with having an open mind, of course, is that people will insist on coming along and trying to put things in it."

-- Terry Pratchett, Diggers

"Time is a drug. Too much of it kills you."

-- Terry Pratchett, Small Gods

"The whole of life is just like watching a film. Only it's as though you always get in ten minutes after the big picture has started, and no-one will tell you the plot, so you have to work it out all yourself from the clues."

-- Terry Pratchett, Moving Pictures
`

// Sample JSON quotes for testing
const SampleJSONQuotes = `[
  {
    "text": "The trouble with having an open mind, of course, is that people will insist on coming along and trying to put things in it.",
    "author": "Terry Pratchett, Diggers"
  },
  {
    "text": "Time is a drug. Too much of it kills you.",
    "author": "Terry Pratchett, Small Gods"
  }
]`

// CreateSampleQuoteFile creates a sample quote file for testing
func CreateSampleQuoteFile(filePath string, format string) error {
	var content string
	switch format {
	case "pqf":
		content = SamplePQFData
	case "json":
		content = SampleJSONQuotes
	default:
		return fmt.Errorf("unknown format: %s", format)
	}
	
	return writeStringToFile(filePath, content)
}

// writeStringToFile writes a string to a file
func writeStringToFile(filePath, content string) error {
	return fmt.Errorf("file operations mocked in tests")
}