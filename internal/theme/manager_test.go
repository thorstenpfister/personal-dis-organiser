package theme

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"personal-disorganizer/internal/testutil"

	"github.com/charmbracelet/lipgloss"
)

func TestNewManager(t *testing.T) {
	tempDir := testutil.TempDir(t)
	
	tests := []struct {
		name        string
		configDir   string
		expectError bool
	}{
		{
			name:        "create manager with valid config dir",
			configDir:   tempDir,
			expectError: false,
		},
		{
			name:        "create manager with empty config dir",
			configDir:   "",
			expectError: false, // Should still work with empty dir
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager, err := NewManager(tt.configDir)
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if !tt.expectError {
				if manager == nil {
					t.Error("Manager should not be nil")
				}
				
				// Should load default Dracula theme
				theme := manager.GetTheme()
				if theme == nil {
					t.Error("Theme should not be nil")
				}
				
				if theme.Name != "dracula" {
					t.Errorf("Expected default theme 'dracula', got '%s'", theme.Name)
				}
				
				// Should have styles created
				styles := manager.GetStyles()
				if styles == nil {
					t.Error("Styles should not be nil")
				}
			}
		})
	}
}

func TestManager_LoadTheme(t *testing.T) {
	tempDir := testutil.TempDir(t)
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	tests := []struct {
		name        string
		themeName   string
		expectError bool
	}{
		{
			name:        "load built-in dracula theme",
			themeName:   "dracula",
			expectError: false,
		},
		{
			name:        "load built-in light theme",
			themeName:   "light",
			expectError: false,
		},
		{
			name:        "load non-existent theme",
			themeName:   "nonexistent",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.LoadTheme(tt.themeName)
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if !tt.expectError {
				theme := manager.GetTheme()
				if theme.Name != tt.themeName {
					t.Errorf("Expected theme name '%s', got '%s'", tt.themeName, theme.Name)
				}
			}
		})
	}
}

func TestManager_LoadCustomTheme(t *testing.T) {
	tempDir := testutil.TempDir(t)
	
	// Create themes directory
	themesDir := filepath.Join(tempDir, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		t.Fatalf("Failed to create themes directory: %v", err)
	}
	
	// Create custom theme file
	customTheme := `{
		"name": "test",
		"background": "#000000",
		"foreground": "#ffffff",
		"primary": "#ff0000",
		"secondary": "#00ff00",
		"accent": "#0000ff",
		"muted": "#888888",
		"success": "#00ff00",
		"warning": "#ffff00",
		"error": "#ff0000"
	}`
	
	customThemePath := filepath.Join(themesDir, "test.json")
	if err := os.WriteFile(customThemePath, []byte(customTheme), 0644); err != nil {
		t.Fatalf("Failed to create custom theme file: %v", err)
	}
	
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	
	// Load custom theme
	err = manager.LoadTheme("test")
	if err != nil {
		t.Errorf("Failed to load custom theme: %v", err)
	}
	
	theme := manager.GetTheme()
	if theme.Name != "test" {
		t.Errorf("Expected theme name 'test', got '%s'", theme.Name)
	}
	
	if theme.Background != "#000000" {
		t.Errorf("Expected background '#000000', got '%s'", theme.Background)
	}
}

func TestManager_GetTheme(t *testing.T) {
	tempDir := testutil.TempDir(t)
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	tests := []struct {
		name        string
		themeName   string
		expectError bool
	}{
		{
			name:        "get dracula theme",
			themeName:   "dracula",
			expectError: false,
		},
		{
			name:        "get light theme",
			themeName:   "light",
			expectError: false,
		},
		{
			name:        "get unknown theme",
			themeName:   "unknown",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme, err := manager.getTheme(tt.themeName)
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if !tt.expectError {
				if theme == nil {
					t.Error("Theme should not be nil")
				}
				
				if theme.Name != tt.themeName {
					t.Errorf("Expected theme name '%s', got '%s'", tt.themeName, theme.Name)
				}
			}
		})
	}
}

func TestManager_GetDraculaTheme(t *testing.T) {
	tempDir := testutil.TempDir(t)
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	theme := manager.getDraculaTheme()
	
	if theme == nil {
		t.Fatal("Dracula theme should not be nil")
	}
	
	if theme.Name != "dracula" {
		t.Errorf("Expected theme name 'dracula', got '%s'", theme.Name)
	}
	
	if theme.Background != "#282a36" {
		t.Errorf("Expected background '#282a36', got '%s'", theme.Background)
	}
	
	if theme.Primary != "#bd93f9" {
		t.Errorf("Expected primary '#bd93f9', got '%s'", theme.Primary)
	}
}

func TestManager_GetLightTheme(t *testing.T) {
	tempDir := testutil.TempDir(t)
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	theme := manager.getLightTheme()
	
	if theme == nil {
		t.Fatal("Light theme should not be nil")
	}
	
	if theme.Name != "light" {
		t.Errorf("Expected theme name 'light', got '%s'", theme.Name)
	}
	
	if theme.Background != "#ffffff" {
		t.Errorf("Expected background '#ffffff', got '%s'", theme.Background)
	}
	
	if theme.Primary != "#007acc" {
		t.Errorf("Expected primary '#007acc', got '%s'", theme.Primary)
	}
}

func TestManager_CreateStyles(t *testing.T) {
	tempDir := testutil.TempDir(t)
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	styles := manager.GetStyles()
	
	if styles == nil {
		t.Fatal("Styles should not be nil")
	}
	
	// Test that all style components exist
	styleTests := []struct {
		name  string
		style lipgloss.Style
	}{
		{"Base", styles.Base},
		{"Title", styles.Title},
		{"StatusBar", styles.StatusBar},
		{"TodayHeader", styles.TodayHeader},
		{"DayHeader", styles.DayHeader},
		{"Secondary", styles.Secondary},
		{"TaskActive", styles.TaskActive},
		{"TaskCompleted", styles.TaskCompleted},
		{"CheckboxActive", styles.CheckboxActive},
		{"CheckboxDone", styles.CheckboxDone},
		{"Calendar", styles.Calendar},
		{"Footer", styles.Footer},
		{"Quote", styles.Quote},
		{"Help", styles.Help},
		{"Search", styles.Search},
	}
	
	for _, tt := range styleTests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic test that style exists (lipgloss.Style is a struct, can't be nil)
			// We can test that it has some properties set
			rendered := tt.style.Render("test")
			if rendered == "" {
				t.Errorf("Style %s should render non-empty string", tt.name)
			}
		})
	}
}

func TestManager_SaveTheme(t *testing.T) {
	tempDir := testutil.TempDir(t)
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Create custom theme
	customTheme := &Theme{
		Name:       "custom_save_test",
		Background: "#123456",
		Foreground: "#abcdef",
		Primary:    "#ff0000",
		Secondary:  "#00ff00",
		Accent:     "#0000ff",
		Muted:      "#888888",
		Success:    "#00cc00",
		Warning:    "#ffcc00",
		Error:      "#cc0000",
	}
	
	// Save theme
	err = manager.SaveTheme(customTheme)
	if err != nil {
		t.Errorf("Failed to save theme: %v", err)
	}
	
	// Verify file was created
	themePath := filepath.Join(tempDir, "themes", "custom_save_test.json")
	if _, err := os.Stat(themePath); os.IsNotExist(err) {
		t.Error("Theme file was not created")
	}
	
	// Verify theme can be loaded back
	err = manager.LoadTheme("custom_save_test")
	if err != nil {
		t.Errorf("Failed to load saved theme: %v", err)
	}
	
	loadedTheme := manager.GetTheme()
	if loadedTheme.Name != customTheme.Name {
		t.Errorf("Expected theme name '%s', got '%s'", customTheme.Name, loadedTheme.Name)
	}
	
	if loadedTheme.Background != customTheme.Background {
		t.Errorf("Expected background '%s', got '%s'", customTheme.Background, loadedTheme.Background)
	}
}

func TestManager_ErrorHandling(t *testing.T) {
	tempDir := testutil.TempDir(t)
	
	// Create themes directory with invalid theme file
	themesDir := filepath.Join(tempDir, "themes")
	if err := os.MkdirAll(themesDir, 0755); err != nil {
		t.Fatalf("Failed to create themes directory: %v", err)
	}
	
	// Create invalid theme file
	invalidThemePath := filepath.Join(themesDir, "invalid.json")
	if err := os.WriteFile(invalidThemePath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to create invalid theme file: %v", err)
	}
	
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}
	
	// Try to load invalid theme - should fall back to built-in
	err = manager.LoadTheme("invalid")
	if err == nil {
		t.Error("Expected error when loading invalid theme")
	}
}

func TestManager_ThemeFileHandling(t *testing.T) {
	tempDir := testutil.TempDir(t)
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	tests := []struct {
		name           string
		setupThemeFile func(string) error
		themeName      string
		expectError    bool
	}{
		{
			name: "load theme from testdata",
			setupThemeFile: func(themesDir string) error {
				// Copy test theme file
				content, err := testutil.ReadTestFile(filepath.Join("testdata", "custom_theme.json"))
				if err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(themesDir, "custom.json"), []byte(content), 0644)
			},
			themeName:   "custom",
			expectError: false,
		},
		{
			name: "load incomplete theme",
			setupThemeFile: func(themesDir string) error {
				content, err := testutil.ReadTestFile(filepath.Join("testdata", "incomplete_theme.json"))
				if err != nil {
					return err
				}
				return os.WriteFile(filepath.Join(themesDir, "incomplete.json"), []byte(content), 0644)
			},
			themeName:   "incomplete",
			expectError: false, // Should load with missing fields as empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create themes directory
			themesDir := filepath.Join(tempDir, "themes")
			if err := os.MkdirAll(themesDir, 0755); err != nil {
				t.Fatalf("Failed to create themes directory: %v", err)
			}
			
			// Setup theme file
			if err := tt.setupThemeFile(themesDir); err != nil {
				t.Skipf("Failed to setup theme file: %v", err)
				return
			}
			
			// Load theme
			err := manager.LoadTheme(tt.themeName)
			
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			
			if !tt.expectError {
				theme := manager.GetTheme()
				if theme.Name != tt.themeName {
					t.Errorf("Expected theme name '%s', got '%s'", tt.themeName, theme.Name)
				}
			}
		})
	}
}

func TestTheme_Validation(t *testing.T) {
	tests := []struct {
		name  string
		theme Theme
	}{
		{
			name: "complete theme",
			theme: Theme{
				Name:       "test",
				Background: "#000000",
				Foreground: "#ffffff",
				Primary:    "#ff0000",
				Secondary:  "#00ff00",
				Accent:     "#0000ff",
				Muted:      "#888888",
				Success:    "#00cc00",
				Warning:    "#ffcc00",
				Error:      "#cc0000",
			},
		},
		{
			name: "theme with empty name",
			theme: Theme{
				Name:       "",
				Background: "#000000",
				Foreground: "#ffffff",
				Primary:    "#ff0000",
				Secondary:  "#00ff00",
				Accent:     "#0000ff",
				Muted:      "#888888",
				Success:    "#00cc00",
				Warning:    "#ffcc00",
				Error:      "#cc0000",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation checks
			if tt.theme.Name == "" && tt.name != "theme with empty name" {
				t.Error("Theme name should not be empty")
			}
			
			// Validate color format (basic check for hex colors)
			colors := []string{
				tt.theme.Background,
				tt.theme.Foreground,
				tt.theme.Primary,
				tt.theme.Secondary,
				tt.theme.Accent,
				tt.theme.Muted,
				tt.theme.Success,
				tt.theme.Warning,
				tt.theme.Error,
			}
			
			for _, color := range colors {
				if color != "" && !strings.HasPrefix(color, "#") {
					t.Errorf("Color should start with #, got: %s", color)
				}
				if color != "" && len(color) != 7 {
					t.Errorf("Color should be 7 characters long, got: %s (%d chars)", color, len(color))
				}
			}
		})
	}
}

func TestManager_StyleConsistency(t *testing.T) {
	tempDir := testutil.TempDir(t)
	manager, err := NewManager(tempDir)
	if err != nil {
		t.Fatalf("Failed to create manager: %v", err)
	}

	// Test that changing themes updates styles
	initialTheme := manager.GetTheme().Name
	initialStyles := manager.GetStyles()
	
	// Load different theme
	newTheme := "light"
	if initialTheme == "light" {
		newTheme = "dracula"
	}
	
	err = manager.LoadTheme(newTheme)
	if err != nil {
		t.Fatalf("Failed to load theme: %v", err)
	}
	
	newStyles := manager.GetStyles()
	
	// Styles should be different (this is a basic check)
	if initialStyles == newStyles {
		t.Error("Styles should be different after theme change")
	}
	
	// Verify theme changed
	if manager.GetTheme().Name == initialTheme {
		t.Error("Theme should have changed")
	}
}