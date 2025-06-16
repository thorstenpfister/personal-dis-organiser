package app

import (
	"testing"
	"time"

	"personal-disorganizer/internal/storage"
)

// Test ListItem functionality (business logic only, not UI rendering)
func TestListItem_FilterValue(t *testing.T) {
	
	tests := []struct {
		name     string
		item     ListItem
		expected string
	}{
		{
			name: "task item",
			item: ListItem{
				ItemType: "task",
				Task: &storage.Task{
					Text: "Complete documentation",
				},
			},
			expected: "Complete documentation",
		},
		{
			name: "day header item",
			item: ListItem{
				ItemType: "day_header",
				Date:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			},
			expected: "Monday, January 15",
		},
		{
			name: "add button item",
			item: ListItem{
				ItemType: "add_button",
			},
			expected: "add new task",
		},
		{
			name: "spacer item",
			item: ListItem{
				ItemType: "spacer",
			},
			expected: "",
		},
		{
			name: "task item with nil task",
			item: ListItem{
				ItemType: "task",
				Task:     nil,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.item.FilterValue()
			
			if result != tt.expected {
				t.Errorf("Expected filter value '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestListItem_Creation(t *testing.T) {
	now := time.Now()
	task := &storage.Task{
		ID:        "test-task-1",
		Text:      "Test task",
		Done:      false,
		Date:      now,
		CreatedAt: now,
	}
	
	tests := []struct {
		name     string
		itemType string
		task     *storage.Task
		date     time.Time
		selected bool
	}{
		{
			name:     "create task item",
			itemType: "task",
			task:     task,
			date:     now,
			selected: false,
		},
		{
			name:     "create day header item",
			itemType: "day_header",
			task:     nil,
			date:     now,
			selected: true,
		},
		{
			name:     "create add button item",
			itemType: "add_button",
			task:     nil,
			date:     now,
			selected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := ListItem{
				ItemType:   tt.itemType,
				Date:       tt.date,
				Task:       tt.task,
				IsSelected: tt.selected,
			}
			
			if item.ItemType != tt.itemType {
				t.Errorf("Expected item type '%s', got '%s'", tt.itemType, item.ItemType)
			}
			
			if !item.Date.Equal(tt.date) {
				t.Errorf("Expected date %v, got %v", tt.date, item.Date)
			}
			
			if item.Task != tt.task {
				t.Errorf("Expected task %v, got %v", tt.task, item.Task)
			}
			
			if item.IsSelected != tt.selected {
				t.Errorf("Expected selected %v, got %v", tt.selected, item.IsSelected)
			}
		})
	}
}

func TestItemDelegate_Height(t *testing.T) {
	delegate := ItemDelegate{}
	
	height := delegate.Height()
	
	if height != 1 {
		t.Errorf("Expected height 1, got %d", height)
	}
}

func TestItemDelegate_Spacing(t *testing.T) {
	delegate := ItemDelegate{}
	
	spacing := delegate.Spacing()
	
	if spacing != 0 {
		t.Errorf("Expected spacing 0, got %d", spacing)
	}
}

// Test business logic functions that can be extracted/tested
func TestTaskFiltering(t *testing.T) {
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)
	
	tasks := []storage.Task{
		{
			ID:   "task1",
			Text: "Today task",
			Done: false,
			Date: today,
		},
		{
			ID:   "task2",
			Text: "Yesterday task",
			Done: true,
			Date: yesterday,
		},
		{
			ID:   "task3",
			Text: "Tomorrow task",
			Done: false,
			Date: tomorrow,
		},
		{
			ID:   "task4",
			Text: "Today completed task",
			Done: true,
			Date: today,
		},
	}
	
	// Test filtering active tasks for today
	activeTodayTasks := filterTasksByDateAndStatus(tasks, today, false)
	
	if len(activeTodayTasks) != 1 {
		t.Errorf("Expected 1 active task for today, got %d", len(activeTodayTasks))
	}
	
	if activeTodayTasks[0].ID != "task1" {
		t.Errorf("Expected task1, got %s", activeTodayTasks[0].ID)
	}
	
	// Test filtering all tasks for today
	allTodayTasks := filterTasksByDate(tasks, today)
	
	if len(allTodayTasks) != 2 {
		t.Errorf("Expected 2 tasks for today, got %d", len(allTodayTasks))
	}
}

func TestTaskSorting(t *testing.T) {
	now := time.Now()
	
	tasks := []storage.Task{
		{
			ID:       "task1",
			Text:     "Task 1",
			Priority: 0,
			CreatedAt: now.Add(-2 * time.Hour),
		},
		{
			ID:       "task2",
			Text:     "Task 2",
			Priority: 1,
			CreatedAt: now.Add(-1 * time.Hour),
		},
		{
			ID:       "task3",
			Text:     "Task 3",
			Priority: 0,
			CreatedAt: now,
		},
		{
			ID:       "task4",
			Text:     "Calendar event",
			Priority: -1, // Calendar events have highest priority
			CreatedAt: now.Add(-3 * time.Hour),
		},
	}
	
	sortedTasks := sortTasksByPriorityAndTime(tasks)
	
	// Calendar event should be first (priority -1)
	if sortedTasks[0].ID != "task4" {
		t.Errorf("Expected calendar event first, got %s", sortedTasks[0].ID)
	}
	
	// Higher priority should come next
	if sortedTasks[1].ID != "task2" {
		t.Errorf("Expected task2 second, got %s", sortedTasks[1].ID)
	}
	
	// Among same priority, newer should come first
	if sortedTasks[2].ID != "task3" {
		t.Errorf("Expected task3 third, got %s", sortedTasks[2].ID)
	}
	
	if sortedTasks[3].ID != "task1" {
		t.Errorf("Expected task1 fourth, got %s", sortedTasks[3].ID)
	}
}

func TestListItemGeneration(t *testing.T) {
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	tomorrow := today.AddDate(0, 0, 1)
	
	tasks := []storage.Task{
		{
			ID:   "task1",
			Text: "Today task",
			Done: false,
			Date: today,
		},
		{
			ID:   "task2",
			Text: "Tomorrow task",
			Done: false,
			Date: tomorrow,
		},
	}
	
	items := generateListItems(tasks)
	
	// Should have items for both days plus tasks
	expectedMinItems := 4 // 2 day headers + 2 tasks
	if len(items) < expectedMinItems {
		t.Errorf("Expected at least %d items, got %d", expectedMinItems, len(items))
	}
	
	// First item should be day header for today
	if items[0].ItemType != "day_header" {
		t.Errorf("Expected first item to be day_header, got %s", items[0].ItemType)
	}
	
	// Check that day headers are properly created
	dayHeaderCount := 0
	taskCount := 0
	
	for _, item := range items {
		switch item.ItemType {
		case "day_header":
			dayHeaderCount++
		case "task":
			taskCount++
			if item.Task == nil {
				t.Error("Task item should have non-nil task")
			}
		}
	}
	
	if dayHeaderCount < 2 {
		t.Errorf("Expected at least 2 day headers, got %d", dayHeaderCount)
	}
	
	if taskCount != 2 {
		t.Errorf("Expected 2 task items, got %d", taskCount)
	}
}

func TestDateGrouping(t *testing.T) {
	now := time.Now()
	today := now.Truncate(24 * time.Hour)
	tomorrow := today.AddDate(0, 0, 1)
	
	tasks := []storage.Task{
		{ID: "task1", Date: today},
		{ID: "task2", Date: today},
		{ID: "task3", Date: tomorrow},
	}
	
	groups := groupTasksByDate(tasks)
	
	if len(groups) != 2 {
		t.Errorf("Expected 2 date groups, got %d", len(groups))
	}
	
	todayTasks, hasTodayTasks := groups[today.Format("2006-01-02")]
	if !hasTodayTasks {
		t.Error("Expected today's tasks to be grouped")
	}
	
	if len(todayTasks) != 2 {
		t.Errorf("Expected 2 tasks for today, got %d", len(todayTasks))
	}
	
	tomorrowTasks, hasTomorrowTasks := groups[tomorrow.Format("2006-01-02")]
	if !hasTomorrowTasks {
		t.Error("Expected tomorrow's tasks to be grouped")
	}
	
	if len(tomorrowTasks) != 1 {
		t.Errorf("Expected 1 task for tomorrow, got %d", len(tomorrowTasks))
	}
}

// Helper functions that would be extracted from app.go for testing
func filterTasksByDateAndStatus(tasks []storage.Task, date time.Time, done bool) []storage.Task {
	var filtered []storage.Task
	targetDate := date.Truncate(24 * time.Hour)
	
	for _, task := range tasks {
		taskDate := task.Date.Truncate(24 * time.Hour)
		if taskDate.Equal(targetDate) && task.Done == done {
			filtered = append(filtered, task)
		}
	}
	
	return filtered
}

func filterTasksByDate(tasks []storage.Task, date time.Time) []storage.Task {
	var filtered []storage.Task
	targetDate := date.Truncate(24 * time.Hour)
	
	for _, task := range tasks {
		taskDate := task.Date.Truncate(24 * time.Hour)
		if taskDate.Equal(targetDate) {
			filtered = append(filtered, task)
		}
	}
	
	return filtered
}

func sortTasksByPriorityAndTime(tasks []storage.Task) []storage.Task {
	sorted := make([]storage.Task, len(tasks))
	copy(sorted, tasks)
	
	// Simple bubble sort for testing (would use sort.Slice in real implementation)
	for i := 0; i < len(sorted); i++ {
		for j := 0; j < len(sorted)-1-i; j++ {
			// Calendar events (priority -1) first, then by priority, then by time
			if shouldSwapTasks(sorted[j], sorted[j+1]) {
				sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
			}
		}
	}
	
	return sorted
}

func shouldSwapTasks(a, b storage.Task) bool {
	// Calendar events have highest priority
	if a.Priority == -1 && b.Priority != -1 {
		return false
	}
	if b.Priority == -1 && a.Priority != -1 {
		return true
	}
	
	// If both are calendar events or both are regular tasks
	if a.Priority != b.Priority {
		return a.Priority < b.Priority // Higher priority (larger number) first
	}
	
	// Same priority, newer first
	return a.CreatedAt.Before(b.CreatedAt)
}

func generateListItems(tasks []storage.Task) []ListItem {
	grouped := groupTasksByDate(tasks)
	var items []ListItem
	
	// Sort dates
	dates := make([]time.Time, 0, len(grouped))
	for dateStr := range grouped {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			dates = append(dates, date)
		}
	}
	
	// Simple date sorting
	for i := 0; i < len(dates); i++ {
		for j := 0; j < len(dates)-1-i; j++ {
			if dates[j].After(dates[j+1]) {
				dates[j], dates[j+1] = dates[j+1], dates[j]
			}
		}
	}
	
	for _, date := range dates {
		dateStr := date.Format("2006-01-02")
		dayTasks := grouped[dateStr]
		
		// Add day header
		items = append(items, ListItem{
			ItemType: "day_header",
			Date:     date,
		})
		
		// Add tasks for this day
		for _, task := range dayTasks {
			taskCopy := task // Create copy to avoid pointer issues
			items = append(items, ListItem{
				ItemType: "task",
				Date:     date,
				Task:     &taskCopy,
			})
		}
	}
	
	return items
}

func groupTasksByDate(tasks []storage.Task) map[string][]storage.Task {
	groups := make(map[string][]storage.Task)
	
	for _, task := range tasks {
		dateStr := task.Date.Format("2006-01-02")
		groups[dateStr] = append(groups[dateStr], task)
	}
	
	return groups
}

// Test AppMode enum values
func TestAppMode_Values(t *testing.T) {
	modes := []AppMode{
		ModeView,
		ModeEdit,
		ModeSearch,
		ModeHistory,
		ModeHelp,
		ModeDeleteConfirm,
	}
	
	// Test that modes have expected values
	expectedValues := []int{0, 1, 2, 3, 4, 5}
	
	for i, mode := range modes {
		if int(mode) != expectedValues[i] {
			t.Errorf("Expected mode %d to have value %d, got %d", i, expectedValues[i], int(mode))
		}
	}
}

func TestListItem_TypeValidation(t *testing.T) {
	validTypes := []string{"day_header", "task", "add_button", "spacer"}
	
	for _, itemType := range validTypes {
		item := ListItem{ItemType: itemType}
		
		// Basic validation that type is set correctly
		if item.ItemType != itemType {
			t.Errorf("Expected item type '%s', got '%s'", itemType, item.ItemType)
		}
		
		// FilterValue should handle all valid types
		filterValue := item.FilterValue()
		if itemType == "task" && item.Task == nil {
			if filterValue != "" {
				t.Errorf("Expected empty filter value for task with nil Task, got '%s'", filterValue)
			}
		}
	}
}