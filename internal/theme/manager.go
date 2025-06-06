package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
)

// Theme represents a color theme configuration
type Theme struct {
	Name       string `json:"name"`
	Background string `json:"background"`
	Foreground string `json:"foreground"`
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Accent     string `json:"accent"`
	Muted      string `json:"muted"`
	Success    string `json:"success"`
	Warning    string `json:"warning"`
	Error      string `json:"error"`
}

// Styles contains all the styled components
type Styles struct {
	Base           lipgloss.Style
	Title          lipgloss.Style
	StatusBar      lipgloss.Style
	TodayHeader    lipgloss.Style
	DayHeader      lipgloss.Style
	Secondary      lipgloss.Style
	TaskActive     lipgloss.Style
	TaskCompleted  lipgloss.Style
	CheckboxActive lipgloss.Style
	CheckboxDone   lipgloss.Style
	Calendar       lipgloss.Style
	Footer         lipgloss.Style
	Quote          lipgloss.Style
	Help           lipgloss.Style
	Search         lipgloss.Style
}

// Manager handles theme loading and style creation
type Manager struct {
	currentTheme *Theme
	styles       *Styles
	configDir    string
}

// NewManager creates a new theme manager
func NewManager(configDir string) (*Manager, error) {
	m := &Manager{
		configDir: configDir,
	}
	
	// Load default Dracula theme
	if err := m.LoadTheme("dracula"); err != nil {
		return nil, fmt.Errorf("failed to load default theme: %w", err)
	}
	
	return m, nil
}

// LoadTheme loads a theme by name
func (m *Manager) LoadTheme(themeName string) error {
	theme, err := m.getTheme(themeName)
	if err != nil {
		return fmt.Errorf("failed to get theme %s: %w", themeName, err)
	}
	
	m.currentTheme = theme
	m.createStyles()
	
	return nil
}

// getTheme retrieves a theme configuration
func (m *Manager) getTheme(themeName string) (*Theme, error) {
	// First try to load from user themes directory
	themePath := filepath.Join(m.configDir, "themes", themeName+".json")
	if data, err := os.ReadFile(themePath); err == nil {
		theme := &Theme{}
		if err := json.Unmarshal(data, theme); err == nil {
			return theme, nil
		}
	}
	
	// Fall back to built-in themes
	switch themeName {
	case "dracula":
		return m.getDraculaTheme(), nil
	case "light":
		return m.getLightTheme(), nil
	default:
		return nil, fmt.Errorf("unknown theme: %s", themeName)
	}
}

// getDraculaTheme returns the built-in Dracula theme
func (m *Manager) getDraculaTheme() *Theme {
	return &Theme{
		Name:       "dracula",
		Background: "#282a36",
		Foreground: "#f8f8f2",
		Primary:    "#bd93f9",
		Secondary:  "#6272a4",
		Accent:     "#ff79c6",
		Muted:      "#44475a",
		Success:    "#50fa7b",
		Warning:    "#f1fa8c",
		Error:      "#ff5555",
	}
}

// getLightTheme returns a built-in light theme
func (m *Manager) getLightTheme() *Theme {
	return &Theme{
		Name:       "light",
		Background: "#ffffff",
		Foreground: "#333333",
		Primary:    "#007acc",
		Secondary:  "#666666",
		Accent:     "#e91e63",
		Muted:      "#cccccc",
		Success:    "#4caf50",
		Warning:    "#ff9800",
		Error:      "#f44336",
	}
}

// createStyles creates all the lipgloss styles based on the current theme
func (m *Manager) createStyles() {
	theme := m.currentTheme
	
	m.styles = &Styles{
		Base: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Background)).
			Foreground(lipgloss.Color(theme.Foreground)),
			
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Primary)).
			Bold(true).
			Padding(0, 1),
			
		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Muted)).
			Foreground(lipgloss.Color(theme.Foreground)).
			Padding(0, 1),
			
		TodayHeader: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Primary)).
			Foreground(lipgloss.Color(theme.Background)).
			Bold(true).
			Padding(0, 1),
			
		DayHeader: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Muted)).
			Foreground(lipgloss.Color(theme.Foreground)).
			Bold(true).
			Padding(0, 1),
			
		Secondary: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Secondary)).
			Bold(true),
			
		TaskActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Foreground)),
			
		TaskCompleted: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Secondary)).
			Strikethrough(true),
			
		CheckboxActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Primary)).
			Bold(true),
			
		CheckboxDone: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Success)).
			Bold(true),
			
		Calendar: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Accent)).
			Italic(true),
			
		Footer: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Muted)).
			Foreground(lipgloss.Color(theme.Foreground)).
			Padding(0, 1).
			Margin(1, 0, 0, 0),
			
		Quote: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Secondary)).
			Italic(true),
			
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Secondary)).
			Italic(true),
			
		Search: lipgloss.NewStyle().
			Background(lipgloss.Color(theme.Accent)).
			Foreground(lipgloss.Color(theme.Background)).
			Bold(true).
			Padding(0, 1),
	}
}

// GetStyles returns the current styles
func (m *Manager) GetStyles() *Styles {
	return m.styles
}

// GetTheme returns the current theme
func (m *Manager) GetTheme() *Theme {
	return m.currentTheme
}

// SaveTheme saves a custom theme to the themes directory
func (m *Manager) SaveTheme(theme *Theme) error {
	themesDir := filepath.Join(m.configDir, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}
	
	themePath := filepath.Join(themesDir, theme.Name+".json")
	data, err := json.MarshalIndent(theme, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal theme: %w", err)
	}
	
	if err := os.WriteFile(themePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme file: %w", err)
	}
	
	return nil
}