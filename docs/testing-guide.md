# Personal Disorganiser - Testing Guide

This guide provides comprehensive documentation for testing the Personal Disorganiser CLI application.

## Overview

The Personal Disorganiser uses a multi-phase testing strategy targeting 80%+ coverage for functional components while maintaining focus on business logic over UI rendering.

## Test Architecture

### Test Structure
```
internal/
├── testutil/           # Shared testing utilities
│   ├── helpers.go      # Test helpers and data generators
│   └── mocks.go        # Mock implementations
├── storage/
│   ├── persistence.go
│   ├── persistence_test.go
│   ├── persistence_bench_test.go
│   └── testdata/       # Sample configuration and data files
├── parser/
│   ├── pratchett.go
│   ├── pratchett_test.go
│   ├── pratchett_bench_test.go
│   └── testdata/       # Sample PQF and JSON quote files
└── [other modules follow similar pattern]
```

### Test Categories

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test component interactions
3. **Benchmark Tests**: Performance testing and optimization
4. **Mock Tests**: External dependency testing

## Running Tests

### Basic Test Commands

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with verbose output
make test-verbose

# Run specific module tests
go test -v ./internal/storage/
go test -v ./internal/parser/
```

### Performance Testing

```bash
# Run all benchmarks
./scripts/run-benchmarks.sh

# Run specific benchmarks
go test -bench=. -benchmem ./internal/search/
go test -bench=BenchmarkEngine_Search ./internal/search/
```

### Coverage Analysis

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# View coverage by function
go tool cover -func=coverage.out
```

## Test Utilities

### Mock Implementations

The `internal/testutil` package provides mock implementations for testing:

```go
// Mock logger for testing error handling
logger := &testutil.MockLogger{}
manager.SetLogger(logger)

// Mock HTTP client for testing calendar integration
client := testutil.NewMockHTTPClient()
client.SetResponse("https://example.com/test.ics", 200, icsContent)

// Mock time provider for deterministic testing
timeProvider := testutil.NewMockTimeProvider(fixedTime)
```

### Test Data Generators

```go
// Generate test tasks
task := testutil.CreateTestTask("id", "text", false, time.Now())

// Generate test configuration
config := testutil.DefaultTestConfig()

// Create temporary directories
tempDir := testutil.TempDir(t)
```

### Assertion Helpers

```go
// File existence checks
testutil.AssertFileExists(t, "/path/to/file")
testutil.AssertFileNotExists(t, "/path/to/file")

// JSON comparison
testutil.AssertJSONEqual(t, expected, actual)
```

## Module-Specific Testing

### Storage Module (80.6% Coverage)

**What's Tested:**
- Configuration loading and saving
- Data persistence (JSON marshalling/unmarshalling)
- Task creation with UUID generation
- Error logging functionality
- Data purge operations

**Key Test Files:**
- `persistence_test.go`: Comprehensive functionality tests
- `persistence_bench_test.go`: Performance benchmarks
- `testdata/`: Sample configuration and data files

### Parser Module (97.4% Coverage)

**What's Tested:**
- PQF format parsing (Terry Pratchett quote files)
- JSON quote loading and validation
- Multi-line quote handling
- Author attribution parsing
- Malformed input graceful handling

**Key Test Files:**
- `pratchett_test.go`: PQF and JSON parsing tests
- `pratchett_bench_test.go`: Parsing performance benchmarks
- `testdata/`: Sample PQF and JSON quote files

### Search Module (100.0% Coverage)

**What's Tested:**
- Fuzzy matching algorithm
- Score calculation and ranking
- Result sorting (by score and date)
- Active task prioritization
- Edge cases and performance

**Key Test Files:**
- `fuzzy_test.go`: Complete search functionality
- `fuzzy_bench_test.go`: Search performance across dataset sizes

### Quotes Module (100.0% Coverage)

**What's Tested:**
- Quote manager initialization
- Multiple file loading (relative/absolute paths)
- Random quote selection
- Error handling for missing/invalid files
- Quote count and availability tracking

### Calendar Module (87.1% Coverage)

**What's Tested:**
- iCal parsing (VEVENT handling)
- Multiple datetime format support
- HTTP integration (with error handling)
- Event filtering by date
- Task conversion logic

**Key Test Files:**
- `ical_test.go`: iCal parsing and HTTP integration
- `testdata/`: Sample iCal files with various formats

### Theme Module (88.6% Coverage)

**What's Tested:**
- Built-in theme loading (Dracula, Light)
- Custom theme file handling
- Style generation with Lipgloss
- Theme validation and error handling
- Theme persistence

**Key Test Files:**
- `manager_test.go`: Theme management functionality
- `testdata/`: Sample theme JSON files

### App Module (Focused Testing)

**What's Tested:**
- ListItem business logic
- Task filtering and sorting algorithms
- Date grouping functionality
- Core data structures and enums

**Note**: UI rendering components are intentionally excluded from testing.

## Best Practices

### Writing Tests

1. **Use Table-Driven Tests**: Organize test cases in structs for clarity
2. **Test Happy Path and Edge Cases**: Cover both successful and error scenarios
3. **Use Meaningful Test Names**: Describe what the test validates
4. **Create Isolated Tests**: Each test should be independent
5. **Use Setup/Teardown**: Properly clean up resources

### Example Test Structure

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name        string
        input       InputType
        expected    OutputType
        expectError bool
    }{
        {
            name:        "successful case",
            input:       validInput,
            expected:    expectedOutput,
            expectError: false,
        },
        {
            name:        "error case",
            input:       invalidInput,
            expected:    OutputType{},
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionName(tt.input)
            
            if tt.expectError && err == nil {
                t.Error("Expected error but got none")
            }
            
            if !tt.expectError && err != nil {
                t.Errorf("Unexpected error: %v", err)
            }
            
            if !tt.expectError && !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

### Performance Testing

1. **Use Realistic Data Sizes**: Test with small, medium, and large datasets
2. **Measure Allocations**: Use `-benchmem` to track memory usage
3. **Reset Timer**: Call `b.ResetTimer()` after setup
4. **Report Allocations**: Call `b.ReportAllocs()` for detailed metrics

```go
func BenchmarkFunction(b *testing.B) {
    // Setup data
    data := generateTestData(1000)
    
    b.ResetTimer()
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        result := Function(data)
        _ = result // Prevent optimization
    }
}
```

## Coverage Goals

| Module | Target | Achieved | Status |
|--------|---------|----------|---------|
| Storage | 80%+ | 80.6% | ✅ |
| Parser | 80%+ | 97.4% | ✅ |
| Search | 80%+ | 100.0% | ✅ |
| Quotes | 80%+ | 100.0% | ✅ |
| Calendar | 80%+ | 87.1% | ✅ |
| Theme | 80%+ | 88.6% | ✅ |
| App Core | Focused | 1.3%* | ✅ |

*App module has low overall coverage due to UI components being excluded from testing strategy.

## Continuous Integration

The project includes GitHub Actions workflow (`.github/workflows/test.yml`) that:

1. Runs all tests with race detection
2. Generates coverage reports
3. Enforces 80% coverage threshold
4. Uploads coverage artifacts
5. Fails builds if coverage drops below threshold

## Troubleshooting

### Common Issues

1. **Import Cycles**: Avoid importing packages under test in testutil
2. **File Permissions**: Ensure test files have proper read/write permissions
3. **Race Conditions**: Use `-race` flag to detect concurrent access issues
4. **Memory Leaks**: Monitor benchmark allocation metrics

### Debug Commands

```bash
# Run tests with race detection
go test -race ./...

# Run tests with timeout
go test -timeout=30s ./...

# Run specific test function
go test -run=TestSpecificFunction ./internal/module/

# Verbose output with test names
go test -v -run=TestPattern ./...
```

## Contributing

When adding new features:

1. Write tests first (TDD approach)
2. Ensure coverage doesn't drop below 80%
3. Add benchmarks for performance-critical code
4. Update test documentation
5. Add sample test data files when needed

## Resources

- [Go Testing Package](https://golang.org/pkg/testing/)
- [Go Benchmark Guide](https://golang.org/pkg/testing/#hdr-Benchmarks)
- [Table Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
- [Test Coverage](https://blog.golang.org/cover)