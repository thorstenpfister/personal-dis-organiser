# Personal Disorganizer - CLI Task Management Tool

## Project Overview

The Personal Disorganizer is a command-line interface task management application designed to provide a simple, keyboard-driven experience for managing daily tasks and calendar events. Despite its playful name, it aims to be a highly organized and efficient productivity tool.

## Core Philosophy

- **Today-focused**: The main view centers around "today" as the primary workspace
- **Keyboard-driven**: All interactions should be possible without a mouse
- **Minimalist**: Clean, distraction-free interface with essential functionality only
- **State-persistent**: Tasks and completion status persist across sessions

## Primary Features

### 1. Today View (Main Interface)
- **Central Focus**: Today's date prominently displayed as the main workspace
- **Calendar Integration**: Import and display events from configurable iCal calendars
- **Task Management**: Simple to-do items that can be added, completed, and managed
- **Unified Display**: Calendar events and tasks shown in a single, organized list

### 3. Task Management System
- **Quick Entry**: Press "+" or Enter on the add button to create new tasks
- **Toggle Completion**: Space bar marks tasks as done/undone with visual feedback
- **Visual States**: 
  - Active tasks: Open checkbox (☐) with normal white text
  - Completed tasks: Checkmark (☑) with gray text and strikethrough
  - Completed tasks fade into background
- **Persistent State**: Task completion status saved between sessions
- **Item Navigation**: Up/Down arrows to navigate between tasks and events within a day

### 4. Navigation & History
- **Item Navigation**: Up/Down arrows navigate between tasks/events within the current day
- **Day Navigation**: Special key combinations or menu options to move between different days
- **History Access**: Two methods to access history:
  - Press Up arrow when at the top of today's list to reveal "show history" option, confirm with Enter
  - Press "h" key for immediate history access
- **Hidden by Default**: Older dates are collapsed to reduce clutter
- **State Preservation**: Historical tasks maintain their completion status

### 4. Calendar Integration
- **iCal Support**: Import multiple calendar sources via configuration file
- **Automatic Sorting**: Calendar events sorted by meeting time at top of list
- **Time Display**: Events show start/end times clearly
- **Configuration-Driven**: Calendar sources managed through central config file

### 6. Fuzzy Search System
- **Search Activation**: Press "/" key to enter fuzzy search mode
- **fzf-like Interface**: Fast, responsive fuzzy search similar to fzf tool
- **Smart Result Ordering**:
  - Active/future tasks displayed at the top of search results
  - Historical/completed tasks shown at the bottom
  - Relevance-based sorting within each section
- **Cross-Date Search**: Search across all dates and tasks in the system
- **Quick Navigation**: Select result and immediately navigate to that task's date
- **Smart Sorting**: 
  - Calendar events at top (sorted by time)
  - Undated to-dos below calendar events
- **Manual Reordering**: Shift + Up/Down arrows to move selected items
- **Priority Management**: Ability to reprioritize tasks within their sections

## Technical Requirements

### Core Technology Stack
- **Language**: Go programming language
- **UI Framework**: Charm's Bubble Tea ecosystem for modern TUI development
  - **bubbletea**: Main TUI framework and runtime
  - **lipgloss**: Styling and layout engine
  - **glamour**: Markdown rendering (for help/documentation)
  - **bubbles**: Pre-built UI components
- **Data Storage**: JSON-based local storage for simplicity and portability
- **Configuration**: JSON configuration files for user customization

### Data Management
- **Local Storage**: All data stored locally in JSON format in ~/.config/personal-disorganizer/
- **State Tracking**: Hierarchical task structure with unlimited nesting levels
- **Configuration Management**: Central config file for user customization
- **Instant Persistence**: Automatic saving on every change to prevent data loss
- **Error Logging**: Timestamped error logs stored in ~/.config/personal-disorganizer/

### File Structure
```
~/.config/personal-disorganizer/
├── config.json          # Main configuration with theme settings
├── data.json            # Task and state data
├── error.log            # Application error log with timestamps
├── quotes/              # Quote files directory
│   └── (user-provided quote files)
└── themes/              # Additional theme definitions (optional)
    ├── dracula.json     # Default Dracula theme
    ├── light.json       # Light theme variant
    └── custom.json      # User custom themes
```

### Application Structure
```
personal-disorganizer/
├── Makefile             # Build, install, uninstall targets
├── cmd/
│   └── main.go          # Main application entry point
├── internal/
│   ├── parser/
│   │   └── pratchett.go # Terry Pratchett quote parser
│   ├── theme/
│   │   └── manager.go   # Theme management and color handling
│   ├── storage/
│   │   └── persistence.go # Data persistence and file management
│   └── ...              # Other internal packages
└── scripts/
    └── fetch-quotes.sh  # Quote fetching script
```

## Detailed Task Breakdown

### Phase 1: Core Infrastructure
1. **Project Setup**
   - Initialize Go module and project structure
   - Set up Bubble Tea framework with bubbletea, lipgloss, glamour, and bubbles
   - Create configuration file loading system with support for multiple quote files
   - Implement JSON data persistence layer
   - Configure lipgloss styles for checkboxes and task states with Dracula theme defaults
   - Create Makefile with build, install, uninstall targets

2. **Build and Installation System**
   - Create comprehensive Makefile with proper targets:
     - `make build`: Compile the application
     - `make install`: Install to system PATH (e.g., `/usr/local/bin`)
     - `make uninstall`: Remove from system PATH
     - `make clean`: Clean build artifacts
     - `make quotes-pratchett`: Download and parse Terry Pratchett quotes
   - Implement cross-platform installation support
   - Create quote file parsing system for Terry Pratchett format

3. **Theme System Foundation**
   - Design theme configuration structure with Dracula as default
   - Create color management system using lipgloss
   - Implement theme validation and fallback mechanisms
   - Set up dynamic styling system for all UI components
   - Create theme color mapping for checkboxes, text states, and UI elements

3. **Basic UI Framework**
   - Create main today view layout using Bubble Tea model-view-update pattern
   - Implement keyboard navigation system with up/down for item selection
   - Set up list display with checkboxes (☐/☑) and proper formatting using lipgloss
   - Apply Dracula theme colors to all UI components
   - Add status bar for current date and context with themed styling
   - Configure visual styling for completed vs incomplete tasks using theme colors

4. **Task Management Foundation**
   - Design task data structure with all required fields
   - Implement task creation, editing, and deletion
   - Create task state management (done/undone toggle with checkboxes)
   - Add visual formatting: ☐ for incomplete, ☑ for complete with strikethrough
   - Apply theme colors to task states (purple checkboxes, green completed, gray faded)
   - Implement item-level navigation within daily view

### Phase 2: Calendar Integration
4. **iCal Parser Implementation**
   - Research and implement iCal file format parsing
   - Create calendar event data structures
   - Implement calendar URL fetching and caching
   - Add error handling for network and parsing issues

5. **Calendar Display Integration**
   - Merge calendar events with task list
   - Implement time-based sorting for events
   - Create clear visual distinction between events and tasks
   - Add proper time formatting and display

### Phase 3: Navigation, History & Search
6. **Date Navigation System**
   - Implement day-by-day navigation (separate from item navigation)
   - Create history view toggle functionality with dual access methods
   - Add date-based task filtering
   - Implement smooth transitions between dates
   - Ensure up/down arrows navigate items within day, not between days

7. **History Management**
   - Create collapsible history view
   - Implement "show history" menu option (Up arrow at top + Enter)
   - Add immediate history access via "h" key
   - Add automatic navigation to historical dates
   - Ensure proper state preservation across date changes

8. **Fuzzy Search Implementation**
   - Create search mode activated by "/" key
   - Implement fzf-like search interface with real-time filtering
   - Design smart result ranking: active/future tasks first, historical last
   - Add cross-date search functionality across all tasks
   - Implement quick navigation from search results to task location

### Phase 4: Advanced Features
9. **Task Reordering**
   - Implement Shift+Arrow key functionality
   - Create drag-and-drop equivalent for terminal
   - Maintain proper sorting rules while allowing manual override
   - Add visual feedback during reordering operations

10. **Quote System Implementation**
    - Create motivational quote display in footer
    - Implement quote rotation on task completion
    - Add support for multiple quote files with random selection
    - Create Terry Pratchett quote parser for the pqf format
    - Implement quote file management and loading system
    - Handle empty quote configuration gracefully (no quotes by default)

### Phase 5: Polish & Optimization
11. **User Experience Enhancements**
    - Add comprehensive keyboard shortcuts
    - Implement proper error handling and user feedback
    - Create help system and documentation using glamour for markdown rendering
    - Add configuration validation and helpful error messages

12. **Performance & Reliability**
    - Optimize data loading and saving operations
    - Implement proper backup and recovery mechanisms
    - Add data validation and corruption prevention
    - Create comprehensive testing suite
    - Optimize fuzzy search performance for large task collections

## Build and Installation System

### Makefile Targets

The project includes a comprehensive Makefile for building, installing, and managing the application:

#### Core Build Targets
- **`make build`**: Compile the application for the current platform
- **`make install`**: Install the binary to system PATH (typically `/usr/local/bin`)
- **`make uninstall`**: Remove the binary from system PATH
- **`make clean`**: Remove build artifacts and temporary files

#### Quote Management Targets
- **`make quotes-pratchett`**: Download and parse Terry Pratchett quotes from lspace.org
- **`make quotes-clean`**: Remove downloaded quote files

#### Development Targets
- **`make test`**: Run the test suite
- **`make fmt`**: Format Go code
- **`make lint`**: Run code linting

### Installation Process

1. **Build**: `make build` compiles the Go application
2. **Install**: `make install` copies the binary to a system PATH location
3. **Usage**: Users can run `personal-disorganizer` from any directory
4. **Uninstall**: `make uninstall` removes the binary from system PATH

### Quote System Architecture

#### Multiple Quote File Support
- **Configuration**: `quote_files` array in config.json specifies quote file paths
- **Random Selection**: Application randomly selects from available quote files
- **No Default Quotes**: No quotes included by default to avoid copyright issues
- **User Control**: Users provide their own quote files or use the Pratchett downloader

#### Terry Pratchett Quote Integration
- **Source**: https://www.lspace.org/ftp/words/pqf/pqf
- **Format Parser**: Custom parser handles the pqf format:
  ```
  Quote text here
  -- Source attribution
  
  
  Next quote text here
  -- Next source attribution
  ```
- **Conversion**: Parser converts to application's JSON quote format:
  ```json
  [
    {
      "text": "Quote text here",
      "author": "Source attribution"
    }
  ]
  ```
- **Make Target**: `make quotes-pratchett` downloads and converts automatically

## Comprehensive Keyboard Shortcuts

### Navigation
- **↑/↓**: Navigate between tasks within current day
- **↑ (at top)**: Show "show history" option (if at very top of today)
- **↓ (at bottom)**: Jump to next day's first task
- **n**: Jump to next day's first task (or "+" if empty)
- **p**: Jump to previous day's first task (or "+" if empty)
- **h**: Jump directly to history view

### Task Management
- **Enter (on "+")**: Create new task and enter edit mode
- **Enter (on task)**: Edit existing task (cursor at end of text)
- **Space**: Toggle task completion (done/undone)
- **Tab**: Indent task one level (create subtask)
- **Shift+Tab**: Outdent task one level (reduce hierarchy)
- **d**: Delete selected task (with confirmation prompt)

### Edit Mode
- **Esc**: Exit edit mode and return to navigation
- **Enter**: Create newline within task text
- **Standard text editing**: Cursor movement, text input, backspace, etc.

### Search and Utility
- **/**: Enter fuzzy search mode (fzf-like interface)
- **Esc (in search)**: Exit search mode
- **Enter (in search)**: Select search result and navigate to task

### Special Navigation
- **Enter (on "show history")**: Expand history view
- **Enter (on delete confirmation)**: Confirm deletion
- **↑/↓ (in confirmation)**: Navigate between "delete" and "cancel"
```json
{
  "calendar_urls": [
    "https://calendar.example.com/user/calendar.ics",
    "webcal://another-calendar.com/feed.ics"
  ],
  "data_file": "data.json",
  "quote_files": [
    "quotes/pratchett.json",
    "quotes/custom.json",
    "quotes/motivational.json"
  ],
  "refresh_interval": 300,
  "date_format": "2006-01-02",
  "time_format": "15:04"
}
```

### data.json Structure
```json
{
  "tasks": [
    {
      "id": "uuid",
      "text": "Task description",
      "done": false,
      "date": "2025-06-06T00:00:00Z",
      "is_calendar": false,
      "start_time": "2025-06-06T09:00:00Z",
      "priority": 0,
      "created_at": "2025-06-06T08:00:00Z"
    }
  ],
  "settings": {
    "last_quote_index": 0,
    "tasks_completed_today": 0
  }
}
```

## UI Layout and Behavior Specifications

### Terminal Layout
- **Minimum Terminal Size**: 80x24 characters
- **Dynamic Proportions**: 
  - Main view: 90% of terminal height
  - Footer: 10% of terminal height (minimum 2 lines)
- **Resize Handling**: Automatic re-render on terminal resize events
- **Text Wrapping**: Configurable width (default 120 characters)
- **Overflow Handling**: Horizontal scroll for content exceeding width

### Visual Hierarchy Display
- **Indentation**: 2 spaces per hierarchy level
- **Maximum Nesting**: Unlimited levels supported
- **Visual Indicators**:
  - Level 0: `☐ Task text`
  - Level 1: `  ☐ Subtask text`
  - Level 2: `    ☐ Sub-subtask text`
- **Completion Display**: `☑` with strikethrough and faded color

### Edit Mode Interface
- **Cursor Display**: Visible cursor with blinking
- **Multi-line Support**: Enter creates newlines within task
- **Text Input**: Full unicode, emoji, and special character support
- **Visual Feedback**: Highlight edit mode with different background color

### Search Interface
- **Search Prompt**: `/` prefix shows search is active
- **Real-time Filtering**: Results update as user types
- **Result Display**: 
  - Active/future tasks at top with normal colors
  - Historical tasks at bottom with muted colors
  - Hierarchy context shown in results

### Error Handling and User Feedback
- **Loading States**: Spinner or indicator for file operations
- **Error Messages**: Clear, actionable error descriptions
- **Confirmation Prompts**: 
  - Delete: "Delete task? [Delete] [Cancel]"
  - Clear highlight on selected option
- **Status Messages**: Brief feedback for successful operations

## Performance and Scalability Specifications

### Data Limits
- **Maximum Tasks**: 10,000 tasks per day (practical limit)
- **Maximum Hierarchy Depth**: 20 levels (reasonable limit)
- **Maximum Task Text**: 2,000 characters per task
- **Historical Data**: 2 years of data (configurable retention)

### Memory Usage
- **Target Memory**: <50MB for typical usage (1000 tasks)
- **Startup Time**: <500ms on modern hardware
- **Response Time**: <100ms for most operations
- **Search Performance**: <200ms for 5000+ tasks

### File Operations
- **Save Frequency**: Immediate on every change
- **Backup Strategy**: Keep last 5 versions of data.json
- **File Locking**: Prevent corruption from multiple instances
- **Error Recovery**: Automatic recovery from corrupted files

## Error Logging Specifications

### Log File Format
```
[2025-06-06 14:30:15] ERROR: Failed to parse config.json: invalid JSON at line 15
[2025-06-06 14:30:16] WARN: Quote file not found: ~/.config/personal-disorganizer/quotes/custom.json
[2025-06-06 14:30:17] INFO: Application started successfully
[2025-06-06 14:30:20] ERROR: Failed to save data.json: permission denied
```

### Error Categories
- **ERROR**: Critical issues that prevent functionality
- **WARN**: Non-critical issues that don't break functionality  
- **INFO**: General application events
- **DEBUG**: Detailed information for troubleshooting (optional)

### Log Management
- **Log Rotation**: Keep last 10 log files, max 1MB each
- **Log Location**: ~/.config/personal-disorganizer/error.log
- **Privacy**: No sensitive user data in logs

## Success Criteria
- [ ] Today view displays current date prominently
- [ ] Tasks show proper checkboxes: ☐ for incomplete, ☑ for completed
- [ ] Hierarchical task structure with unlimited nesting via Tab/Shift+Tab
- [ ] Tasks can be added, completed, and organized hierarchically
- [ ] Up/Down arrows navigate between items within a day
- [ ] Navigation between dates works smoothly with boundary detection
- [ ] History view toggles properly (Up at top + Enter, or "h" key)
- [ ] Fuzzy search works via "/" with fzf-like interface
- [ ] Search results prioritize active/future tasks over historical ones
- [ ] Application installs system-wide via Makefile
- [ ] Quote system supports multiple user-provided quote files
- [ ] Terry Pratchett quotes can be downloaded and parsed via make target
- [ ] Dracula theme applied as default with full customization support
- [ ] All UI elements respect theme configuration
- [ ] Files stored in ~/.config/personal-disorganizer/ directory
- [ ] Dynamic layout responds to terminal resize (90% main, 10% footer)
- [ ] Configurable text wrapping (default 120 characters)
- [ ] Error logging with timestamps
- [ ] Full emoji, unicode, and newline support in tasks
- [ ] Instant persistence on every change
- [ ] All keyboard shortcuts function as specified
- [ ] Data persists correctly between sessions for incomplete, ☑ for completed
- [ ] Tasks can be added, completed, and reordered
- [ ] Up/Down arrows navigate between items within a day
- [ ] Navigation between dates works smoothly (separate from item navigation)
- [ ] History view toggles properly (Up at top + Enter, or "h" key)
- [ ] Fuzzy search works via "/" with fzf-like interface
- [ ] Search results prioritize active/future tasks over historical ones
- [ ] Application installs system-wide via Makefile
- [ ] Quote system supports multiple user-provided quote files
- [ ] Terry Pratchett quotes can be downloaded and parsed via make target
- [ ] Dracula theme applied as default with full customization support
- [ ] All UI elements respect theme configuration
- [ ] All keyboard shortcuts function as specified
- [ ] Data persists correctly between sessions

### User Experience Goals
- [ ] Interface is intuitive and requires minimal learning
- [ ] All operations feel fast and responsive
- [ ] Visual feedback is clear and helpful
- [ ] Error handling is graceful and informative
- [ ] Application starts quickly and reliably

### Technical Standards
- [ ] Code is well-documented and maintainable
- [ ] Configuration system is flexible and robust
- [ ] Data storage is reliable and recoverable
- [ ] Memory usage is reasonable for a CLI tool
- [ ] Cross-platform compatibility (Windows, macOS, Linux)

## Implementation Details and Missing Specifications

### Dependencies and Libraries
- **Go Version**: Go 1.19+ required
- **External Dependencies**:
  - `github.com/charmbracelet/bubbletea` - TUI framework
  - `github.com/charmbracelet/lipgloss` - Styling
  - `github.com/charmbracelet/glamour` - Markdown rendering
  - `github.com/charmbracelet/bubbles` - UI components
  - `github.com/google/uuid` - UUID generation
- **Standard Library Usage**: 
  - `encoding/json` - Data persistence
  - `os` - File operations
  - `time` - Date/time handling
  - `strings` - Text processing

### Cross-Platform Considerations
- **Config Directory Creation**: Handle different OS permissions
- **Path Handling**: Use `filepath.Join()` for cross-platform paths
- **Terminal Capabilities**: Graceful degradation for limited terminals
- **Signal Handling**: Proper cleanup on SIGINT/SIGTERM

### Security Considerations
- **File Permissions**: 
  - Config files: 0644 (readable by user and group)
  - Data files: 0644 (readable by user and group)
  - Log files: 0644 (readable by user and group)
- **Input Validation**: Sanitize all user input to prevent issues
- **Path Traversal**: Validate file paths to prevent directory traversal
- **Resource Limits**: Prevent memory exhaustion from large inputs

### Testing Strategy
- **Unit Tests**: Core logic and data structures
- **Integration Tests**: File operations and persistence
- **UI Tests**: Keyboard navigation and display logic
- **Performance Tests**: Large dataset handling
- **Cross-Platform Tests**: Verify functionality across OS

### Development Workflow
- **Code Organization**: Clear separation of concerns
- **Documentation**: Comprehensive code comments
- **Version Control**: Semantic versioning for releases
- **Build Process**: Automated builds and testing
- **Release Process**: Tagged releases with binaries

## Future Enhancement Opportunities

- **Additional Themes**: More built-in themes (Solarized, Nord, GitHub, etc.)
- **Export**: Export capabilities for tasks and completed work logs
- **Categories**: Task categorization and filtering system
- **Advanced Search**: Search with filters (by date range, completion status, categories)
- **Search History**: Remember recent searches and provide quick access