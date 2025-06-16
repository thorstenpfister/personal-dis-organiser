package search

import (
	"testing"
	"time"

	"personal-disorganizer/internal/storage"
	"personal-disorganizer/internal/testutil"
)

// Benchmark search performance with various dataset sizes
func BenchmarkEngine_Search_Small(b *testing.B) {
	benchmarkSearchWithTaskCount(b, 10)
}

func BenchmarkEngine_Search_Medium(b *testing.B) {
	benchmarkSearchWithTaskCount(b, 100)
}

func BenchmarkEngine_Search_Large(b *testing.B) {
	benchmarkSearchWithTaskCount(b, 1000)
}

func BenchmarkEngine_Search_ExtraLarge(b *testing.B) {
	benchmarkSearchWithTaskCount(b, 10000)
}

func benchmarkSearchWithTaskCount(b *testing.B, taskCount int) {
	engine := NewEngine()
	tasks := generateBenchmarkTasks(taskCount)
	query := "project"
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		results := engine.Search(query, tasks)
		_ = results // Avoid optimization
	}
}

// Benchmark calculate score performance
func BenchmarkEngine_CalculateScore_ShortText(b *testing.B) {
	engine := NewEngine()
	query := "test"
	text := "test task"
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		score := engine.calculateScore(query, text)
		_ = score
	}
}

func BenchmarkEngine_CalculateScore_LongText(b *testing.B) {
	engine := NewEngine()
	query := "project"
	text := "This is a very long project description that contains multiple words and should test the performance of the fuzzy matching algorithm when dealing with lengthy text content that might be found in real-world applications"
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		score := engine.calculateScore(query, text)
		_ = score
	}
}

// Benchmark search with different query lengths
func BenchmarkEngine_Search_ShortQuery(b *testing.B) {
	benchmarkSearchWithQuery(b, "t")
}

func BenchmarkEngine_Search_MediumQuery(b *testing.B) {
	benchmarkSearchWithQuery(b, "project")
}

func BenchmarkEngine_Search_LongQuery(b *testing.B) {
	benchmarkSearchWithQuery(b, "long project description with multiple words")
}

func benchmarkSearchWithQuery(b *testing.B, query string) {
	engine := NewEngine()
	tasks := generateBenchmarkTasks(1000)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		results := engine.Search(query, tasks)
		_ = results
	}
}

// Generate tasks for benchmarking
func generateBenchmarkTasks(count int) []storage.Task {
	tasks := make([]storage.Task, count)
	now := time.Now()
	
	taskTexts := []string{
		"Complete project documentation",
		"Review code changes",
		"Meeting with team members",
		"Update project timeline and milestones",
		"Write comprehensive unit tests",
		"Fix authentication bug in login system",
		"Deploy application to staging environment",
		"Create detailed user manual and guides",
		"Optimize database performance queries",
		"Implement new feature for user management",
		"Refactor legacy code components",
		"Design system architecture diagrams",
		"Conduct security audit and testing",
		"Backup database and configuration files",
		"Monitor application performance metrics",
	}
	
	for i := 0; i < count; i++ {
		tasks[i] = storage.Task{
			ID:        testutil.MockUUID(i),
			Text:      taskTexts[i%len(taskTexts)],
			Done:      i%4 == 0, // 25% completion rate
			Date:      now.AddDate(0, 0, i%30-15), // Spread across month
			Priority:  i % 3,
			CreatedAt: now.Add(time.Duration(-i) * time.Hour),
			Level:     0,
		}
	}
	
	return tasks
}