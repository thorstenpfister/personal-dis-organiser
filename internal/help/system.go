package help

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

// System handles help documentation
type System struct {
	renderer *glamour.TermRenderer
}

// NewSystem creates a new help system
func NewSystem() (*System, error) {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return nil, err
	}
	
	return &System{
		renderer: renderer,
	}, nil
}

// GetHelpText returns formatted help documentation
func (h *System) GetHelpText() (string, error) {
	markdown := `# Personal Disorganizer - Help

## Navigation
- **↑/↓ or k/j**: Navigate between tasks within current day
- **n**: Go to next day
- **p**: Go to previous day
- **h**: View history of all tasks

## Task Management
- **Enter**: Edit selected task or add new task (when on "+")
- **Space**: Toggle task completion (☐ ↔ ☑)
- **d**: Delete selected task
- **Tab**: Indent task (increase hierarchy level)
- **Shift+Tab**: Outdent task (decrease hierarchy level)

## Task Reordering
- **Shift+↑**: Move task up (within day or to previous day)
- **Shift+↓**: Move task down (within day or to next day)
- Cross-day movement: Tasks moved beyond day boundaries transfer to adjacent days
- Boundary: Cannot move tasks to dates before today

## Search
- **/**: Enter search mode
- In search mode:
  - Type to search across all tasks
  - **↑/↓**: Navigate search results
  - **Enter**: Go to selected task
  - **Esc**: Exit search

## Edit Mode
- **Enter**: Save changes
- **Esc**: Cancel editing
- Standard text editing (cursor movement, backspace, etc.)

## Quotes
- **r**: Refresh quote (get new random quote)

## Other
- **q or Ctrl+C**: Quit application

## Configuration

The application stores data in **~/.config/personal-disorganizer/**:
- **config.json**: Main configuration
- **data.json**: Task and completion data
- **quotes/**: Optional quote files
- **themes/**: Custom theme definitions

## Calendar Integration

Add calendar URLs to config.json:
` + "```json" + `
{
  "calendar_urls": [
    "https://calendar.example.com/feed.ics",
    "webcal://another-calendar.com/feed.ics"
  ]
}
` + "```" + `

## Quote System

To add Terry Pratchett quotes, run:
` + "```bash" + `
make quotes-pratchett
` + "```" + `

Or add your own quote files to the quotes/ directory.

## Themes

The default theme is Dracula. Create custom themes in the themes/ directory:
` + "```json" + `
{
  "name": "custom",
  "background": "#1e1e1e",
  "foreground": "#d4d4d4",
  "primary": "#569cd6",
  "accent": "#c586c0"
}
` + "```" + `

Press **q** to close help.`

	return h.renderer.Render(markdown)
}

// GetKeyboardShortcuts returns a condensed list of shortcuts
func (h *System) GetKeyboardShortcuts() string {
	return strings.Join([]string{
		"Navigation: ↑/↓/n/p/h",
		"Tasks: Enter/Space/d/Tab",
		"Reorder: Shift+↑/↓",
		"Search: /",
		"Quit: q",
	}, " • ")
}