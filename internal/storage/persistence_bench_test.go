package storage

import (
	"os"
	"testing"
	"time"

	"personal-disorganizer/internal/testutil"
)

// Benchmark storage operations
func BenchmarkStorage_SaveData_Small(b *testing.B) {
	benchmarkSaveData(b, 10)
}

func BenchmarkStorage_SaveData_Medium(b *testing.B) {
	benchmarkSaveData(b, 100)
}

func BenchmarkStorage_SaveData_Large(b *testing.B) {
	benchmarkSaveData(b, 1000)
}

func benchmarkSaveData(b *testing.B, taskCount int) {
	tempDir := testutil.TempDir(&testing.T{}) // Use empty testing.T for benchmark
	dataPath := tempDir + "/data.json"
	
	storage := &Storage{
		configDir: tempDir,
		dataPath:  dataPath,
	}
	
	// Generate test data
	data := generateBenchmarkAppData(taskCount)
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		err := storage.SaveData(data)
		if err != nil {
			b.Fatalf("SaveData failed: %v", err)
		}
		
		// Clean up for next iteration
		os.Remove(dataPath)
	}
}

func BenchmarkStorage_LoadData_Small(b *testing.B) {
	benchmarkLoadData(b, 10)
}

func BenchmarkStorage_LoadData_Medium(b *testing.B) {
	benchmarkLoadData(b, 100)
}

func BenchmarkStorage_LoadData_Large(b *testing.B) {
	benchmarkLoadData(b, 1000)
}

func benchmarkLoadData(b *testing.B, taskCount int) {
	tempDir := testutil.TempDir(&testing.T{})
	dataPath := tempDir + "/data.json"
	
	storage := &Storage{
		configDir: tempDir,
		dataPath:  dataPath,
	}
	
	// Pre-create test data file
	data := generateBenchmarkAppData(taskCount)
	err := storage.SaveData(data)
	if err != nil {
		b.Fatalf("Failed to setup test data: %v", err)
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		loadedData, err := storage.LoadData()
		if err != nil {
			b.Fatalf("LoadData failed: %v", err)
		}
		_ = loadedData
	}
}

// Benchmark task creation
func BenchmarkStorage_CreateTask(b *testing.B) {
	storage := &Storage{}
	taskText := "Benchmark task creation performance"
	taskDate := time.Now()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		task := storage.CreateTask(taskText, taskDate)
		_ = task
	}
}

// Generate app data for benchmarking
func generateBenchmarkAppData(taskCount int) *AppData {
	tasks := make([]Task, taskCount)
	now := time.Now()
	
	taskTexts := []string{
		"Complete project documentation",
		"Review code changes",
		"Meeting with team",
		"Update project timeline",
		"Write unit tests",
		"Fix authentication bug",
		"Deploy to staging",
		"Create user manual",
		"Optimize database queries",
		"Implement new features",
	}
	
	for i := 0; i < taskCount; i++ {
		tasks[i] = Task{
			ID:        testutil.MockUUID(i),
			Text:      taskTexts[i%len(taskTexts)],
			Done:      i%4 == 0,
			Date:      now.AddDate(0, 0, i%30-15),
			IsCalendar: i%10 == 0, // 10% calendar events
			Priority:  i % 3,
			CreatedAt: now.Add(time.Duration(-i) * time.Hour),
			Level:     i % 3, // 0-2 hierarchy levels
		}
		
		if tasks[i].IsCalendar {
			tasks[i].StartTime = now.Add(time.Duration(i) * time.Hour)
			tasks[i].Priority = -1
		}
	}
	
	return &AppData{
		Tasks: tasks,
		Settings: Settings{
			LastQuoteIndex:      taskCount % 100,
			TasksCompletedToday: taskCount / 10,
		},
	}
}