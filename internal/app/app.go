package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"personal-disorganizer/internal/calendar"
	"personal-disorganizer/internal/help"
	"personal-disorganizer/internal/parser"
	"personal-disorganizer/internal/quotes"
	"personal-disorganizer/internal/search"
	"personal-disorganizer/internal/storage"
	"personal-disorganizer/internal/theme"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AppMode represents the current mode of the application
type AppMode int

const (
	ModeView AppMode = iota
	ModeEdit
	ModeSearch
	ModeHistory
	ModeHelp
	ModeDeleteConfirm
)

// ListItem represents an item in the list (either a task or a day header)
type ListItem struct {
	ItemType   string        // "day_header", "task", "add_button", "spacer"
	Date       time.Time     // The date this item belongs to
	Task       *storage.Task // The task (nil for day headers and add buttons)
	IsSelected bool          // Whether this item is currently selected
}

// FilterValue implements list.Item interface
func (i ListItem) FilterValue() string {
	switch i.ItemType {
	case "task":
		if i.Task != nil {
			return i.Task.Text
		}
	case "day_header":
		return i.Date.Format("Monday, January 2")
	case "add_button":
		return "add new task"
	}
	return ""
}

// ItemDelegate handles rendering of list items
type ItemDelegate struct {
	styles *theme.Styles
	width  int
}

// Height implements list.ItemDelegate interface
func (d ItemDelegate) Height() int {
	return 1
}

// Spacing implements list.ItemDelegate interface  
func (d ItemDelegate) Spacing() int {
	return 0
}

// Update implements list.ItemDelegate interface
func (d ItemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

// Render implements list.ItemDelegate interface
func (d ItemDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	listItem, ok := item.(ListItem)
	if !ok {
		return
	}
	
	isSelected := index == m.Index()
	
	switch listItem.ItemType {
	case "day_header":
		d.renderDayHeader(w, listItem, isSelected)
	case "task":
		d.renderTask(w, listItem, isSelected)
	case "add_button":
		d.renderAddButton(w, listItem, isSelected)
	}
}

func (d ItemDelegate) renderDayHeader(w io.Writer, item ListItem, selected bool) {
	dateHeader := item.Date.Format("Monday, January 2")
	isToday := item.Date.Truncate(24*time.Hour).Equal(time.Now().Truncate(24*time.Hour))
	isTomorrow := item.Date.Truncate(24*time.Hour).Equal(time.Now().Add(24*time.Hour).Truncate(24*time.Hour))
	
	if isToday {
		dateHeader = "Today - " + dateHeader
		fmt.Fprint(w, d.styles.TodayHeader.Width(d.width).Render(dateHeader))
	} else {
		if isTomorrow {
			dateHeader = "Tomorrow - " + dateHeader
		}
		fmt.Fprint(w, d.styles.DayHeader.Width(d.width).Render(dateHeader))
	}
}

func (d ItemDelegate) renderTask(w io.Writer, item ListItem, selected bool) {
	if item.Task == nil {
		return
	}
	
	task := *item.Task
	prefix := "  "
	if selected {
		prefix = "> "
	}
	
	// Indentation for hierarchy
	indent := strings.Repeat("  ", task.Level)
	
	// Handle calendar events differently
	if task.IsCalendar {
		timeStr := task.StartTime.Format("15:04")
		text := d.styles.Calendar.Render(fmt.Sprintf("%s %s", timeStr, task.Text))
		fmt.Fprintf(w, "%s%s📅 %s", prefix, indent, text)
		return
	}
	
	// Regular task checkbox
	var checkbox string
	var textStyle lipgloss.Style
	
	if task.Done {
		checkbox = d.styles.CheckboxDone.Render("☑")
		textStyle = d.styles.TaskCompleted
	} else {
		checkbox = d.styles.CheckboxActive.Render("☐")
		textStyle = d.styles.TaskActive
	}
	
	text := textStyle.Render(task.Text)
	fmt.Fprintf(w, "%s%s%s %s", prefix, indent, checkbox, text)
}

func (d ItemDelegate) renderAddButton(w io.Writer, item ListItem, selected bool) {
	addButton := "+ Add new task"
	if selected {
		addButton = d.styles.CheckboxActive.Render("> ") + addButton
	} else {
		addButton = "  " + addButton
	}
	fmt.Fprint(w, addButton)
}

// Model represents the entire application state
type Model struct {
	// Core state
	mode       AppMode
	width      int
	height     int
	
	// Data
	storage       *storage.Storage
	appData       *storage.AppData
	tasks         []storage.Task
	calendarTasks []storage.Task
	
	// Managers
	themeManager    *theme.Manager
	quoteManager    *quotes.Manager
	calendarManager *calendar.Manager
	searchEngine    *search.Engine
	helpSystem      *help.System
	
	// UI
	styles    *theme.Styles
	textInput textinput.Model
	list      list.Model
	delegate  ItemDelegate
	
	// View state
	currentDate   time.Time
	showHistory   bool
	searchQuery   string
	searchResults []search.Result
	searchCursor  int
	
	// Edit state
	editDate        time.Time
	editTaskForDate *storage.Task
	
	// Delete confirmation state
	deleteTaskID string
	
	// Quote state
	currentQuote *parser.Quote
	
	// Error handling
	lastError string
}

// NewModel creates a new application model
func NewModel() (*Model, error) {
	// Initialize storage
	storage, err := storage.NewStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}
	
	// Load application data
	appData, err := storage.LoadData()
	if err != nil {
		return nil, fmt.Errorf("failed to load data: %w", err)
	}
	
	// Get config directory
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "personal-disorganizer")
	
	// Initialize theme manager
	themeManager, err := theme.NewManager(configDir)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize theme: %w", err)
	}
	
	// Initialize quote manager
	config := storage.GetConfig()
	quoteManager, err := quotes.NewManager(configDir, config.QuoteFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize quotes: %w", err)
	}
	
	// Initialize calendar manager
	calendarManager := calendar.NewManager(config.CalendarURLs)
	calendarManager.SetLogger(storage)
	
	// Initialize search engine
	searchEngine := search.NewEngine()
	
	// Initialize help system
	helpSystem, err := help.NewSystem()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize help system: %w", err)
	}
	
	// Create text input for editing
	ti := textinput.New()
	ti.Placeholder = "Enter task..."
	
	// Create list component
	delegate := ItemDelegate{styles: themeManager.GetStyles(), width: 80} // Default width
	taskList := list.New([]list.Item{}, delegate, 0, 0)
	taskList.Title = ""
	taskList.SetShowStatusBar(false)
	taskList.SetFilteringEnabled(false)
	taskList.SetShowHelp(false)
	taskList.SetShowTitle(false)
	
	m := &Model{
		mode:            ModeView,
		storage:         storage,
		appData:         appData,
		themeManager:    themeManager,
		quoteManager:    quoteManager,
		calendarManager: calendarManager,
		searchEngine:    searchEngine,
		helpSystem:      helpSystem,
		styles:          themeManager.GetStyles(),
		textInput:       ti,
		list:            taskList,
		delegate:        delegate,
		currentDate:     time.Now().Truncate(24 * time.Hour),
		showHistory:     false,
	}
	
	// Initialize quote if available
	if quoteManager.HasQuotes() {
		m.currentQuote = quoteManager.GetRandomQuote()
	}
	
	m.updateTasksForCurrentDate()
	m.rebuildListItems()
	
	return m, nil
}

// Init initializes the application
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textInput.Width = msg.Width - 4
		m.list.SetWidth(msg.Width)
		
		// Update the delegate width for proper day header rendering
		m.delegate.width = msg.Width
		m.list.SetDelegate(m.delegate)
		
		// Update list height based on current footer size
		m.updateListHeight()
		
	case tea.KeyMsg:
		// Handle text input first if in edit mode (except for special keys)
		if m.mode == ModeEdit {
			switch msg.String() {
			case "esc", "enter":
				// Let these be handled by the mode handler
				return m.handleKeyMsg(msg)
			default:
				// Pass all other keys to text input
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		}
		return m.handleKeyMsg(msg)
	}
	
	return m, nil
}

// handleKeyMsg handles keyboard input
func (m *Model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case ModeView:
		return m.handleViewMode(msg)
	case ModeEdit:
		return m.handleEditMode(msg)
	case ModeSearch:
		return m.handleSearchMode(msg)
	case ModeHistory:
		return m.handleHistoryMode(msg)
	case ModeHelp:
		return m.handleHelpMode(msg)
	case ModeDeleteConfirm:
		return m.handleDeleteConfirmMode(msg)
	}
	return m, nil
}

// handleViewMode handles input in view mode
func (m *Model) handleViewMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
		
	case "enter":
		// Handle enter based on selected list item
		selectedItem := m.getSelectedListItem()
		if selectedItem == nil {
			break
		}
		
		switch selectedItem.ItemType {
		case "add_button":
			m.startEditingNewTaskForDate(selectedItem.Date)
		case "task":
			if selectedItem.Task != nil {
				m.startEditingExistingTask(selectedItem.Task, selectedItem.Date)
			}
		}
		
	case " ":
		// Toggle task completion
		selectedItem := m.getSelectedListItem()
		if selectedItem != nil && selectedItem.ItemType == "task" && selectedItem.Task != nil {
			m.toggleTaskById(selectedItem.Task.ID)
			m.saveData()
			m.rebuildListItemsPreservingSelection()
		}
		
	case "d":
		// Delete task - show confirmation
		selectedItem := m.getSelectedListItem()
		if selectedItem != nil && selectedItem.ItemType == "task" && selectedItem.Task != nil {
			m.deleteTaskID = selectedItem.Task.ID
			m.mode = ModeDeleteConfirm
		}
		
	case "tab":
		// Indent task (increase hierarchy level)
		selectedItem := m.getSelectedListItem()
		if selectedItem != nil && selectedItem.ItemType == "task" && selectedItem.Task != nil {
			m.adjustTaskLevel(selectedItem.Task.ID, 1)
			m.saveData()
			m.rebuildListItemsPreservingSelection()
		}
		
	case "shift+tab":
		// Outdent task (decrease hierarchy level)
		selectedItem := m.getSelectedListItem()
		if selectedItem != nil && selectedItem.ItemType == "task" && selectedItem.Task != nil {
			m.adjustTaskLevel(selectedItem.Task.ID, -1)
			m.saveData()
			m.rebuildListItemsPreservingSelection()
		}
		
	case "shift+up":
		// Move task up (possibly to previous day)
		selectedItem := m.getSelectedListItem()
		if selectedItem != nil && selectedItem.ItemType == "task" && selectedItem.Task != nil && !selectedItem.Task.IsCalendar {
			m.moveTaskUp(selectedItem.Date, selectedItem.Task.ID)
			m.saveData()
			m.rebuildListItemsPreservingSelection()
		}
		
	case "shift+down":
		// Move task down (possibly to next day)
		selectedItem := m.getSelectedListItem()
		if selectedItem != nil && selectedItem.ItemType == "task" && selectedItem.Task != nil && !selectedItem.Task.IsCalendar {
			m.moveTaskDown(selectedItem.Date, selectedItem.Task.ID)
			m.saveData()
			m.rebuildListItemsPreservingSelection()
		}
		
	case "h":
		// Jump to history
		m.mode = ModeHistory
		
	case "?":
		// Show help
		m.mode = ModeHelp
		
	case "r":
		// Refresh quote manually
		m.refreshQuote()
		
	case "/":
		// Enter search mode
		m.mode = ModeSearch
		m.textInput.SetValue("")
		m.textInput.Focus()
		
	case "n":
		// Next day
		m.currentDate = m.currentDate.Add(24 * time.Hour)
		m.updateTasksForCurrentDate()
		m.rebuildListItems()
		
	case "p":
		// Previous day
		m.currentDate = m.currentDate.Add(-24 * time.Hour)
		m.updateTasksForCurrentDate()
		m.rebuildListItems()
		
	default:
		// Let the list handle navigation (up/down/etc)
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	
	return m, nil
}

// handleEditMode handles input in edit mode
func (m *Model) handleEditMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeView
		m.textInput.Blur()
		
	case "enter":
		text := strings.TrimSpace(m.textInput.Value())
		if text != "" {
			if m.editTaskForDate == nil {
				// Creating new task - use smart insertion to preserve hierarchy
				task := m.storage.CreateTask(text, m.editDate)
				m.insertTaskAtPosition(task, m.editDate)
			} else {
				// Editing existing task
				for i := range m.appData.Tasks {
					if m.appData.Tasks[i].ID == m.editTaskForDate.ID {
						m.appData.Tasks[i].Text = text
						break
					}
				}
			}
			m.saveData()
			m.updateTasksForCurrentDate()
			m.rebuildListItems()
		}
		m.mode = ModeView
		m.textInput.Blur()
		m.textInput.SetValue("")
		m.editTaskForDate = nil
	}
	
	return m, nil
}

// handleSearchMode handles input in search mode
func (m *Model) handleSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = ModeView
		m.textInput.Blur()
		m.searchQuery = ""
		m.searchResults = []search.Result{}
		m.searchCursor = 0
		
	case "enter":
		// Navigate to selected search result
		if m.searchCursor < len(m.searchResults) {
			result := m.searchResults[m.searchCursor]
			
			// Find the task in the list and set cursor to it
			// The list should always start from today, not change
			m.setListCursorToTask(result.Task.ID)
			
			m.mode = ModeView
			m.textInput.Blur()
			m.searchQuery = ""
			m.searchResults = []search.Result{}
			m.searchCursor = 0
		}
		
	case "up", "k":
		if m.searchCursor > 0 {
			m.searchCursor--
		}
		
	case "down", "j":
		if m.searchCursor < len(m.searchResults)-1 {
			m.searchCursor++
		}
		
	default:
		// Update search query and results
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(tea.KeyMsg(msg))
		m.searchQuery = m.textInput.Value()
		m.updateSearchResults()
		return m, cmd
	}
	
	return m, nil
}

// handleHistoryMode handles input in history mode
func (m *Model) handleHistoryMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "h":
		m.mode = ModeView
	}
	
	return m, nil
}

// handleHelpMode handles input in help mode
func (m *Model) handleHelpMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "?":
		m.mode = ModeView
	}
	
	return m, nil
}

// handleDeleteConfirmMode handles input in delete confirmation mode
func (m *Model) handleDeleteConfirmMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Confirm deletion
		if m.deleteTaskID != "" {
			m.deleteTaskById(m.deleteTaskID)
			m.saveData()
			m.rebuildListItemsPreservingSelection()
		}
		m.deleteTaskID = ""
		m.mode = ModeView
		
	case "n", "N", "esc":
		// Cancel deletion
		m.deleteTaskID = ""
		m.mode = ModeView
	}
	
	return m, nil
}

// startEditingNewTaskForDate starts editing a new task for a specific date
func (m *Model) startEditingNewTaskForDate(date time.Time) {
	m.mode = ModeEdit
	m.editTaskForDate = nil
	m.editDate = date
	m.textInput.SetValue("")
	m.textInput.Focus()
}

// insertTaskAtPosition inserts a new task at the appropriate position preserving hierarchy
func (m *Model) insertTaskAtPosition(newTask *storage.Task, targetDate time.Time) {
	// Get the currently selected item to determine insertion context
	selectedItem := m.getSelectedListItem()
	
	// Get all tasks for the target date, sorted by priority
	dayTasks := m.getTasksForDate(targetDate)
	
	if selectedItem == nil || len(dayTasks) == 0 {
		// No selection or no existing tasks - just add with default priority
		newTask.Priority = 1
		m.appData.Tasks = append(m.appData.Tasks, *newTask)
		return
	}
	
	if selectedItem.ItemType == "add_button" {
		// Adding at the end of the day - set priority lower than the lowest existing task
		minPriority := 0
		for _, task := range dayTasks {
			if !task.IsCalendar && task.Priority < minPriority {
				minPriority = task.Priority
			}
		}
		newTask.Priority = minPriority - 1
		m.appData.Tasks = append(m.appData.Tasks, *newTask)
		return
	}
	
	if selectedItem.ItemType == "task" && selectedItem.Task != nil {
		// Insert after the selected task and its entire subtask block
		selectedTask := selectedItem.Task
		
		// Find the end of the selected task's subtask block
		// Subtasks have higher level numbers and lower priority numbers (appear immediately after)
		endPriority := selectedTask.Priority - 1
		
		// Look for any subtasks (children) of the selected task
		for _, task := range dayTasks {
			if !task.IsCalendar && 
			   task.Priority < selectedTask.Priority && 
			   task.Level > selectedTask.Level {
				// This is a subtask - find the lowest priority among all subtasks
				if task.Priority < endPriority {
					endPriority = task.Priority
				}
			}
		}
		
		// Insert the new task after the entire block (selected task + its subtasks)
		newTask.Priority = endPriority - 1
		m.appData.Tasks = append(m.appData.Tasks, *newTask)
		return
	}
	
	// Fallback: add at end
	newTask.Priority = 1
	m.appData.Tasks = append(m.appData.Tasks, *newTask)
}


// updateListHeight recalculates and sets the list height based on current footer size
func (m *Model) updateListHeight() {
	if m.width == 0 || m.height == 0 {
		return
	}
	
	// Calculate current footer height
	footer := m.renderFooter()
	footerLines := strings.Count(footer, "\n") + 1
	padding := 2 // content padding + footer padding
	
	// Calculate available height for list
	availableHeight := m.height - footerLines - padding
	if availableHeight < 1 {
		availableHeight = 1
	}
	
	// Update the list height
	m.list.SetHeight(availableHeight)
}

// refreshQuote gets a new random quote
func (m *Model) refreshQuote() {
	if m.quoteManager.HasQuotes() {
		m.currentQuote = m.quoteManager.GetRandomQuote()
		// Update list height since footer size may have changed
		m.updateListHeight()
	}
}

// deleteTask deletes a task
func (m *Model) deleteTask(index int) {
	if index < len(m.tasks) {
		taskID := m.tasks[index].ID
		for i := range m.appData.Tasks {
			if m.appData.Tasks[i].ID == taskID {
				m.appData.Tasks = append(m.appData.Tasks[:i], m.appData.Tasks[i+1:]...)
				break
			}
		}
		m.updateTasksForCurrentDate()
	}
}

// moveTaskUp moves a task up within its day or to the previous day
func (m *Model) moveTaskUp(date time.Time, taskID string) {
	tasks := m.getTasksForDate(date)
	localIndex := -1
	
	// Find the task's current position
	for i, task := range tasks {
		if task.ID == taskID {
			localIndex = i
			break
		}
	}
	
	if localIndex == -1 || tasks[localIndex].IsCalendar {
		return
	}
	
	if localIndex > 0 {
		// Move within the same day
		m.moveTaskWithinDay(taskID, date, localIndex, localIndex-1)
	} else {
		// Move to previous day (only if not moving to past)
		prevDate := date.Add(-24 * time.Hour)
		if !prevDate.Before(m.currentDate) {
			m.moveTaskToDay(taskID, date, prevDate, -1) // -1 means to the end
		}
	}
}

// moveTaskDown moves a task down within its day or to the next day
func (m *Model) moveTaskDown(date time.Time, taskID string) {
	tasks := m.getTasksForDate(date)
	localIndex := -1
	
	// Find the task's current position
	for i, task := range tasks {
		if task.ID == taskID {
			localIndex = i
			break
		}
	}
	
	if localIndex == -1 || tasks[localIndex].IsCalendar {
		return
	}
	
	if localIndex < len(tasks)-1 {
		// Move within the same day
		m.moveTaskWithinDay(taskID, date, localIndex, localIndex+1)
	} else {
		// Move to next day
		nextDate := date.Add(24 * time.Hour)
		m.moveTaskToDay(taskID, date, nextDate, 0) // 0 means to the beginning
	}
}

// moveTaskWithinDay moves a task to a different position within the same day
func (m *Model) moveTaskWithinDay(taskID string, date time.Time, fromIndex, toIndex int) {
	for i := range m.appData.Tasks {
		if m.appData.Tasks[i].ID == taskID {
			// Adjust priority based on direction
			if toIndex < fromIndex {
				m.appData.Tasks[i].Priority++
			} else {
				m.appData.Tasks[i].Priority--
			}
			break
		}
	}
	m.updateTasksForCurrentDate()
}

// moveTaskToDay moves a task from one day to another
func (m *Model) moveTaskToDay(taskID string, fromDate, toDate time.Time, position int) {
	for i := range m.appData.Tasks {
		if m.appData.Tasks[i].ID == taskID {
			// Change the task's date
			m.appData.Tasks[i].Date = toDate
			
			// Set priority based on position
			if position == -1 {
				// Move to end of day - find lowest priority for that day
				minPriority := 0
				for _, task := range m.appData.Tasks {
					if task.Date.Truncate(24*time.Hour).Equal(toDate.Truncate(24*time.Hour)) && !task.IsCalendar {
						if task.Priority < minPriority {
							minPriority = task.Priority
						}
					}
				}
				m.appData.Tasks[i].Priority = minPriority - 1
			} else {
				// Move to beginning of day - find highest priority for that day
				maxPriority := 0
				for _, task := range m.appData.Tasks {
					if task.Date.Truncate(24*time.Hour).Equal(toDate.Truncate(24*time.Hour)) && !task.IsCalendar {
						if task.Priority > maxPriority {
							maxPriority = task.Priority
						}
					}
				}
				m.appData.Tasks[i].Priority = maxPriority + 1
			}
			break
		}
	}
	m.updateTasksForCurrentDate()
}


// updateTasksForCurrentDate filters tasks for the current date
func (m *Model) updateTasksForCurrentDate() {
	m.tasks = []storage.Task{}
	
	// Add calendar events for the current date
	if calendarTasks, err := m.calendarManager.FetchEvents(m.currentDate); err == nil {
		m.calendarTasks = calendarTasks
		m.tasks = append(m.tasks, calendarTasks...)
	}
	
	// Add regular tasks for the current date
	for _, task := range m.appData.Tasks {
		if task.Date.Truncate(24*time.Hour).Equal(m.currentDate) {
			m.tasks = append(m.tasks, task)
		}
	}
	
	// Sort tasks: calendar events first (by time), then regular tasks
	sort.Slice(m.tasks, func(i, j int) bool {
		if m.tasks[i].IsCalendar != m.tasks[j].IsCalendar {
			return m.tasks[i].IsCalendar // Calendar events first
		}
		if m.tasks[i].IsCalendar && m.tasks[j].IsCalendar {
			return m.tasks[i].StartTime.Before(m.tasks[j].StartTime)
		}
		return m.tasks[i].Priority > m.tasks[j].Priority
	})
}

// updateSearchResults updates the search results based on current query
func (m *Model) updateSearchResults() {
	if m.searchQuery == "" {
		m.searchResults = []search.Result{}
		return
	}
	
	m.searchResults = m.searchEngine.Search(m.searchQuery, m.appData.Tasks)
	m.searchCursor = 0
}

// saveData saves the current application data
func (m *Model) saveData() {
	if err := m.storage.SaveData(m.appData); err != nil {
		m.lastError = err.Error()
		m.storage.LogError(err)
	}
}

// View renders the application UI
func (m *Model) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading..."
	}
	
	// Calculate available space
	footer := m.renderFooter()
	
	// Count lines in footer
	footerLines := strings.Count(footer, "\n") + 1
	
	// Reserve space for padding between sections
	padding := 2 // content padding + footer padding
	
	// Calculate available height for main content
	availableHeight := m.height - footerLines - padding
	if availableHeight < 1 {
		availableHeight = 1
	}
	
	var b strings.Builder
	
	// Main content
	switch m.mode {
	case ModeEdit:
		content := m.renderEditView()
		content = m.fitContentToHeight(content, availableHeight)
		b.WriteString(content)
		
		// Add spacing to push footer to bottom
		contentLines := strings.Count(content, "\n") + 1
		remainingLines := availableHeight - contentLines
		if remainingLines > 0 {
			b.WriteString(strings.Repeat("\n", remainingLines))
		}
	case ModeSearch:
		content := m.renderSearchView()
		content = m.fitContentToHeight(content, availableHeight)
		b.WriteString(content)
		
		// Add spacing to push footer to bottom
		contentLines := strings.Count(content, "\n") + 1
		remainingLines := availableHeight - contentLines
		if remainingLines > 0 {
			b.WriteString(strings.Repeat("\n", remainingLines))
		}
	case ModeHistory:
		content := m.renderHistoryView()
		content = m.fitContentToHeight(content, availableHeight)
		b.WriteString(content)
		
		// Add spacing to push footer to bottom
		contentLines := strings.Count(content, "\n") + 1
		remainingLines := availableHeight - contentLines
		if remainingLines > 0 {
			b.WriteString(strings.Repeat("\n", remainingLines))
		}
	case ModeHelp:
		content := m.renderHelpView()
		content = m.fitContentToHeight(content, availableHeight)
		b.WriteString(content)
		
		// Add spacing to push footer to bottom
		contentLines := strings.Count(content, "\n") + 1
		remainingLines := availableHeight - contentLines
		if remainingLines > 0 {
			b.WriteString(strings.Repeat("\n", remainingLines))
		}
	case ModeDeleteConfirm:
		content := m.renderDeleteConfirmView()
		content = m.fitContentToHeight(content, availableHeight)
		b.WriteString(content)
		
		// Add spacing to push footer to bottom
		contentLines := strings.Count(content, "\n") + 1
		remainingLines := availableHeight - contentLines
		if remainingLines > 0 {
			b.WriteString(strings.Repeat("\n", remainingLines))
		}
	default:
		// Use the list component for the main view
		b.WriteString(m.list.View())
	}
	
	// Footer at bottom
	b.WriteString("\n")
	b.WriteString(footer)
	
	return b.String()
}

// fitContentToHeight ensures content fits within the available height
func (m *Model) fitContentToHeight(content string, maxHeight int) string {
	lines := strings.Split(content, "\n")
	
	if len(lines) <= maxHeight {
		return content
	}
	
	// If content is too long, truncate and add indicator
	truncated := lines[:maxHeight-1]
	moreLines := len(lines) - len(truncated)
	
	truncated = append(truncated, fmt.Sprintf("... (%d more lines - use scroll or resize terminal)", moreLines))
	
	return strings.Join(truncated, "\n")
}




// getTasksForDate gets all tasks for a specific date
func (m *Model) getTasksForDate(date time.Time) []storage.Task {
	var tasks []storage.Task
	targetDate := date.Truncate(24 * time.Hour)
	
	for _, task := range m.appData.Tasks {
		if task.Date.Truncate(24*time.Hour).Equal(targetDate) {
			tasks = append(tasks, task)
		}
	}
	
	// Sort tasks: calendar events first (by time), then regular tasks
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].IsCalendar != tasks[j].IsCalendar {
			return tasks[i].IsCalendar // Calendar events first
		}
		if tasks[i].IsCalendar && tasks[j].IsCalendar {
			return tasks[i].StartTime.Before(tasks[j].StartTime)
		}
		return tasks[i].Priority > tasks[j].Priority
	})
	
	return tasks
}



// renderEditView renders the edit mode view
func (m *Model) renderEditView() string {
	var b strings.Builder
	
	if m.editTaskForDate == nil {
		b.WriteString("Adding new task:\n\n")
	} else {
		b.WriteString("Editing task:\n\n")
	}
	
	b.WriteString(m.textInput.View())
	b.WriteString("\n\nPress Enter to save, Esc to cancel")
	
	return b.String()
}

// renderSearchView renders the search mode view
func (m *Model) renderSearchView() string {
	var b strings.Builder
	
	b.WriteString("Search: ")
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")
	
	if len(m.searchResults) == 0 {
		if m.searchQuery != "" {
			b.WriteString("No results found")
		} else {
			b.WriteString("Type to search...")
		}
	} else {
		b.WriteString(fmt.Sprintf("Found %d results:\n\n", len(m.searchResults)))
		
		for i, result := range m.searchResults {
			prefix := " "
			if i == m.searchCursor {
				prefix = ">"
			}
			
			dateStr := result.Task.Date.Format("2006-01-02")
			status := "☐"
			if result.Task.Done {
				status = "☑"
			}
			
			line := fmt.Sprintf("%s %s %s [%s]", prefix, status, result.Match, dateStr)
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	
	b.WriteString("\n↑/↓: navigate • Enter: go to task • Esc: cancel")
	
	return b.String()
}

// renderHistoryView renders the history mode view
func (m *Model) renderHistoryView() string {
	var b strings.Builder
	
	b.WriteString("Task History\n\n")
	
	// Group tasks by date
	tasksByDate := make(map[string][]storage.Task)
	for _, task := range m.appData.Tasks {
		dateKey := task.Date.Format("2006-01-02")
		tasksByDate[dateKey] = append(tasksByDate[dateKey], task)
	}
	
	// Sort dates in reverse order
	var dates []string
	for date := range tasksByDate {
		dates = append(dates, date)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	
	// Display recent dates
	maxDates := 10
	for i, date := range dates {
		if i >= maxDates {
			break
		}
		
		b.WriteString(fmt.Sprintf("=== %s ===\n", date))
		
		tasks := tasksByDate[date]
		for _, task := range tasks {
			status := "☐"
			if task.Done {
				status = "☑"
			}
			b.WriteString(fmt.Sprintf("  %s %s\n", status, task.Text))
		}
		b.WriteString("\n")
	}
	
	b.WriteString("Press Esc or 'h' to return")
	
	return b.String()
}

// renderHelpView renders the help mode view
func (m *Model) renderHelpView() string {
	helpText, err := m.helpSystem.GetHelpText()
	if err != nil {
		return "Error loading help: " + err.Error()
	}
	return helpText
}

// renderDeleteConfirmView renders the delete confirmation view
func (m *Model) renderDeleteConfirmView() string {
	var b strings.Builder
	
	// Find the task being deleted
	var taskText string
	for _, task := range m.appData.Tasks {
		if task.ID == m.deleteTaskID {
			taskText = task.Text
			break
		}
	}
	
	b.WriteString("Delete Task\n\n")
	if taskText != "" {
		b.WriteString(fmt.Sprintf("Are you sure you want to delete this task?\n\n\"%s\"\n\n", taskText))
	} else {
		b.WriteString("Are you sure you want to delete this task?\n\n")
	}
	b.WriteString("Press 'y' to confirm, 'n' or Esc to cancel")
	
	return b.String()
}

// renderFooter renders the application footer
func (m *Model) renderFooter() string {
	var b strings.Builder
	
	// Help text first - make it adaptive to terminal width
	help := "↑/↓: navigate • Shift+↑/↓: move tasks • Enter: edit • Space: toggle • d: delete • h: history • /: search • r: quote • ?: help • q: quit"
	
	// If terminal is narrow, use shorter help text
	if m.width < 130 {
		help = "↑/↓: nav • Shift+↑/↓: move • Enter: edit • Space: toggle • d: del • h: hist • /: search • r: quote • ?: help • q: quit"
	}
	if m.width < 110 {
		help = "↑/↓: nav • Enter: edit • Space: toggle • d: del • h: hist • /: search • r: quote • ?: help • q: quit"
	}
	if m.width < 90 {
		help = "↑/↓/Enter/Space/d/h/r/?/q - Press ? for help"
	}
	
	if m.lastError != "" {
		help = "Error: " + m.lastError
	}
	
	b.WriteString(m.styles.Footer.Width(m.width).Render(help))
	
	// Quote below help interface (if available)
	if m.currentQuote != nil {
		b.WriteString("\n\n") // Visual spacing between help and quote
		b.WriteString(m.renderQuote())
	}
	
	return b.String()
}

// renderQuote renders a properly formatted, centered quote
func (m *Model) renderQuote() string {
	if m.currentQuote == nil {
		return ""
	}
	
	var b strings.Builder
	
	// Format quote text (without author initially)
	quoteText := fmt.Sprintf("\"%s\"", m.currentQuote.Text)
	
	// Use 2/3 of screen width for line breaking, with minimum margin
	maxLineWidth := (m.width * 2) / 3
	if maxLineWidth < 40 {
		maxLineWidth = m.width - 8 // Fallback for very narrow terminals
	}
	
	wrappedQuote := m.wrapText(quoteText, maxLineWidth)
	quoteLines := strings.Split(wrappedQuote, "\n")
	
	// Limit to maximum 4 lines for quote text
	maxQuoteLines := 4
	if len(quoteLines) > maxQuoteLines {
		quoteLines = quoteLines[:maxQuoteLines-1]
		lastLine := quoteLines[len(quoteLines)-1]
		if len(lastLine) > maxLineWidth-3 {
			lastLine = lastLine[:maxLineWidth-3]
		}
		quoteLines[len(quoteLines)-1] = lastLine + "..."
	}
	
	// Center each line of the quote
	for _, line := range quoteLines {
		centeredLine := m.centerText(line, m.width)
		b.WriteString(m.styles.Quote.Render(centeredLine))
		b.WriteString("\n")
	}
	
	// Add author on separate line, also centered
	authorLine := fmt.Sprintf("— %s", m.currentQuote.Author)
	centeredAuthor := m.centerText(authorLine, m.width)
	b.WriteString(m.styles.Quote.Render(centeredAuthor))
	
	return b.String()
}

// centerText centers text within the given width
func (m *Model) centerText(text string, width int) string {
	textLen := len(text)
	if textLen >= width {
		return text
	}
	
	padding := (width - textLen) / 2
	return strings.Repeat(" ", padding) + text
}

// wrapText wraps text to fit within the specified width
func (m *Model) wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	
	var lines []string
	var currentLine strings.Builder
	
	for _, word := range words {
		// Check if adding this word would exceed the width
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " "
		}
		testLine += word
		
		if len(testLine) <= width {
			// Word fits, add it to current line
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		} else {
			// Word doesn't fit, start new line
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
			currentLine.WriteString(word)
		}
	}
	
	// Add the last line if it has content
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	
	return strings.Join(lines, "\n")
}

// rebuildListItems creates the list items for the next 30 days starting from today
func (m *Model) rebuildListItems() {
	var items []list.Item
	
	// Always start from the actual current date (today), not m.currentDate
	today := time.Now().Truncate(24 * time.Hour)
	
	// Add current day (today)
	items = append(items, ListItem{
		ItemType: "day_header",
		Date:     today,
	})
	
	// Add today's tasks (use m.currentDate for task filtering to maintain compatibility)
	todayTasks := m.getTasksForDate(today)
	for _, task := range todayTasks {
		items = append(items, ListItem{
			ItemType: "task",
			Date:     today,
			Task:     &task,
		})
	}
	
	// Add today's "add task" button
	items = append(items, ListItem{
		ItemType: "add_button",
		Date:     today,
	})
	
	// Add next 30 days
	for i := 1; i <= 30; i++ {
		futureDate := today.Add(time.Duration(i) * 24 * time.Hour)
		futureTasks := m.getTasksForDate(futureDate)
		
		// Add day header
		items = append(items, ListItem{
			ItemType: "day_header",
			Date:     futureDate,
		})
		
		// Add tasks for this day
		for _, task := range futureTasks {
			items = append(items, ListItem{
				ItemType: "task",
				Date:     futureDate,
				Task:     &task,
			})
		}
		
		// Add "add task" button for this day
		items = append(items, ListItem{
			ItemType: "add_button",
			Date:     futureDate,
		})
	}
	
	m.list.SetItems(items)
}

// getSelectedListItem returns the currently selected list item
func (m *Model) getSelectedListItem() *ListItem {
	selectedIndex := m.list.Index()
	items := m.list.Items()
	
	if selectedIndex >= 0 && selectedIndex < len(items) {
		if item, ok := items[selectedIndex].(ListItem); ok {
			return &item
		}
	}
	return nil
}

// startEditingExistingTask starts editing an existing task
func (m *Model) startEditingExistingTask(task *storage.Task, date time.Time) {
	m.mode = ModeEdit
	m.editTaskForDate = task
	m.editDate = date
	m.textInput.SetValue(task.Text)
	m.textInput.Focus()
}

// toggleTaskById toggles a task's completion status by ID
func (m *Model) toggleTaskById(taskID string) {
	for i := range m.appData.Tasks {
		if m.appData.Tasks[i].ID == taskID {
			m.appData.Tasks[i].Done = !m.appData.Tasks[i].Done
			// Get a new quote when task status changes
			m.refreshQuote()
			break
		}
	}
	m.updateTasksForCurrentDate()
}

// deleteTaskById deletes a task by ID
func (m *Model) deleteTaskById(taskID string) {
	for i := range m.appData.Tasks {
		if m.appData.Tasks[i].ID == taskID {
			m.appData.Tasks = append(m.appData.Tasks[:i], m.appData.Tasks[i+1:]...)
			break
		}
	}
	m.updateTasksForCurrentDate()
}

// adjustTaskLevel adjusts a task's hierarchy level by ID
func (m *Model) adjustTaskLevel(taskID string, delta int) {
	for i := range m.appData.Tasks {
		if m.appData.Tasks[i].ID == taskID {
			newLevel := m.appData.Tasks[i].Level + delta
			if newLevel >= 0 {
				m.appData.Tasks[i].Level = newLevel
			}
			break
		}
	}
	m.updateTasksForCurrentDate()
}

// setListCursorToTask finds a task by ID in the list and sets the cursor to it
func (m *Model) setListCursorToTask(taskID string) {
	items := m.list.Items()
	for i, item := range items {
		if listItem, ok := item.(ListItem); ok {
			if listItem.ItemType == "task" && listItem.Task != nil && listItem.Task.ID == taskID {
				m.list.Select(i)
				return
			}
		}
	}
}

// rebuildListItemsPreservingSelection rebuilds the list while trying to preserve the current selection
func (m *Model) rebuildListItemsPreservingSelection() {
	// Store the currently selected item info
	var selectedTaskID string
	var selectedDate time.Time
	var selectedItemType string
	
	selectedItem := m.getSelectedListItem()
	if selectedItem != nil {
		selectedItemType = selectedItem.ItemType
		selectedDate = selectedItem.Date
		if selectedItem.Task != nil {
			selectedTaskID = selectedItem.Task.ID
		}
	}
	
	// Rebuild the list
	m.rebuildListItems()
	
	// Try to restore selection
	if selectedTaskID != "" {
		m.setListCursorToTask(selectedTaskID)
	} else if selectedItemType == "add_button" {
		// Find the add button for the same date
		items := m.list.Items()
		for i, item := range items {
			if listItem, ok := item.(ListItem); ok {
				if listItem.ItemType == "add_button" && listItem.Date.Equal(selectedDate) {
					m.list.Select(i)
					return
				}
			}
		}
	} else if selectedItemType == "day_header" {
		// Find the day header for the same date
		items := m.list.Items()
		for i, item := range items {
			if listItem, ok := item.(ListItem); ok {
				if listItem.ItemType == "day_header" && listItem.Date.Equal(selectedDate) {
					m.list.Select(i)
					return
				}
			}
		}
	}
}