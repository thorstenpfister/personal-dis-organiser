# Personal Disorganiser

A modern, keyboard-driven CLI task management application built with Go and Bubble Tea.

![](promo/demo.gif)

## Features

- **Today-focused Interface**: Main view centers around today as your primary workspace
- **Hierarchical Tasks**: Unlimited nesting levels with Tab/Shift+Tab indentation and smart block preservation
- **Calendar Integration**: Import iCal calendars and display events alongside tasks
- **Fuzzy Search**: Fast, fzf-like search across all tasks and dates
- **Task Management**: Create, edit, delete, and reorder tasks with intuitive keyboard shortcuts
- **Quote System**: Optional motivational quotes with Terry Pratchett integration
- **Dracula Theme**: Beautiful default theme with full customization support
- **Data Persistence**: Automatic saving with JSON-based local storage

## Installation

```bash
# Build the application
make build

# Install system-wide
make install

# Optional: Add Terry Pratchett quotes
make quotes-pratchett
```

## Usage

```bash
# Run the application
personal-disorganiser
```

### Keyboard Shortcuts

- **Navigation**: ↑/↓ (navigate tasks), n/p (next/previous day), h (history)
- **Tasks**: Enter (edit), Space (toggle done), d (delete), Tab (indent)
- **Reordering**: Shift+↑/↓ (move tasks up/down)
- **Search**: / (enter search mode)
- **Help**: ? (show comprehensive help)
- **Quit**: q or Ctrl+C

## Configuration

Configuration files are stored in `~/.config/personal-disorganizer/`:

- `config.json` - Main configuration
- `data.json` - Tasks and state data
- `quotes/` - Quote files directory
- `themes/` - Custom theme definitions

### Adding Calendar Integration

Edit `~/.config/personal-disorganizer/config.json`:

```json
{
  "calendar_urls": [
    "https://calendar.example.com/feed.ics",
    "webcal://another-calendar.com/feed.ics"
  ],
  "quote_files": [
    "quotes/pratchett.json",
    "quotes/custom.json"
  ]
}
```

### Custom Themes

Create theme files in `~/.config/personal-disorganizer/themes/`:

```json
{
  "name": "custom",
  "background": "#1e1e1e",
  "foreground": "#d4d4d4",
  "primary": "#569cd6",
  "secondary": "#6a9955",
  "accent": "#c586c0",
  "success": "#4fc1ff",
  "warning": "#ffcc02",
  "error": "#f44747"
}
```

## Build System

- `make build` - Compile the application
- `make install` - Install to system PATH and configure shell
- `make uninstall` - Remove from system PATH
- `make path-check` - Check if install directory is in PATH
- `make source-path` - Show commands to update PATH in current session
- `make clean` - Clean build artifacts
- `make quotes-pratchett` - Download and parse Terry Pratchett quotes
- `make quotes-clean` - Remove quote files
- `make test` - Run tests
- `make fmt` - Format code
- `make lint` - Run linter

## Architecture

```
personal-disorganiser/
├── cmd/main.go              # Application entry point
├── internal/
│   ├── app/app.go          # Main application logic
│   ├── storage/            # Data persistence
│   ├── theme/              # Theme management
│   ├── calendar/           # iCal integration
│   ├── search/             # Fuzzy search
│   ├── quotes/             # Quote system
│   ├── help/               # Help system
│   └── parser/             # Quote file parsing
├── Makefile                # Build automation
├── tasks.md                # Current development tasks
└── docs/                   # Documentation and archives
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [Glamour](https://github.com/charmbracelet/glamour) - Markdown rendering
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components
- [UUID](https://github.com/google/uuid) - Unique identifiers

## License

See LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## Philosophy

Despite its playful name, Personal Disorganiser aims to be a highly organized and efficient productivity tool that:

- Focuses on today while maintaining historical context
- Provides keyboard-driven efficiency without mouse dependency
- Maintains a clean, distraction-free interface
- Persists state reliably across sessions
- Integrates seamlessly with existing calendar workflows