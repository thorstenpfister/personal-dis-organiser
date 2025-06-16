package parser

import (
	"os"
	"strings"
	"testing"
)

// Benchmark PQF parsing performance
func BenchmarkParsePQF_Small(b *testing.B) {
	benchmarkParsePQF(b, 10)
}

func BenchmarkParsePQF_Medium(b *testing.B) {
	benchmarkParsePQF(b, 100)
}

func BenchmarkParsePQF_Large(b *testing.B) {
	benchmarkParsePQF(b, 1000)
}

func benchmarkParsePQF(b *testing.B, quoteCount int) {
	// Generate test PQF content
	content := generatePQFContent(quoteCount)
	
	// Write to temporary file
	tmpFile, err := os.CreateTemp("", "benchmark_*.pqf")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	if _, err := tmpFile.WriteString(content); err != nil {
		b.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		quotes, err := ParsePQF(tmpFile.Name())
		if err != nil {
			b.Fatalf("ParsePQF failed: %v", err)
		}
		_ = quotes
	}
}

// Benchmark JSON quote loading
func BenchmarkLoadQuotes_Small(b *testing.B) {
	benchmarkLoadQuotes(b, 10)
}

func BenchmarkLoadQuotes_Medium(b *testing.B) {
	benchmarkLoadQuotes(b, 100)
}

func BenchmarkLoadQuotes_Large(b *testing.B) {
	benchmarkLoadQuotes(b, 1000)
}

func benchmarkLoadQuotes(b *testing.B, quoteCount int) {
	// Generate test JSON content
	content := generateJSONQuotes(quoteCount)
	
	// Write to temporary file
	tmpFile, err := os.CreateTemp("", "benchmark_*.json")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()
	
	if _, err := tmpFile.WriteString(content); err != nil {
		b.Fatalf("Failed to write temp file: %v", err)
	}
	tmpFile.Close()
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		quotes, err := LoadQuotes(tmpFile.Name())
		if err != nil {
			b.Fatalf("LoadQuotes failed: %v", err)
		}
		_ = quotes
	}
}

// Generate PQF content for benchmarking
func generatePQFContent(quoteCount int) string {
	var builder strings.Builder
	
	quotes := []string{
		"The trouble with having an open mind, of course, is that people will insist on coming along and trying to put things in it.",
		"Time is a drug. Too much of it kills you.",
		"In the beginning there was nothing, which exploded.",
		"Five exclamation marks, the sure sign of an insane mind.",
		"The whole of life is just like watching a film.",
		"Fantasy is the impossible made probable.",
		"Real stupidity beats artificial intelligence every time.",
		"A good bookshop is just a genteel Black Hole that knows how to read.",
		"The pen is mightier than the sword if the sword is very short, and the pen is very sharp.",
		"Give a man a fire and he's warm for a day, but set fire to him and he's warm for the rest of his life.",
	}
	
	authors := []string{
		"Terry Pratchett, Diggers",
		"Terry Pratchett, Small Gods",
		"Terry Pratchett, Lords and Ladies",
		"Terry Pratchett, Reaper Man",
		"Terry Pratchett, Moving Pictures",
		"Terry Pratchett, The Color of Magic",
		"Terry Pratchett, Hogfather",
		"Terry Pratchett, Good Omens",
		"Terry Pratchett, Guards! Guards!",
		"Terry Pratchett, Night Watch",
	}
	
	for i := 0; i < quoteCount; i++ {
		quote := quotes[i%len(quotes)]
		author := authors[i%len(authors)]
		
		builder.WriteString(`"`)
		builder.WriteString(quote)
		builder.WriteString(`"`)
		builder.WriteString("\n\n-- ")
		builder.WriteString(author)
		builder.WriteString("\n\n")
	}
	
	return builder.String()
}

// Generate JSON quotes for benchmarking
func generateJSONQuotes(quoteCount int) string {
	var builder strings.Builder
	
	quotes := []string{
		"The trouble with having an open mind, of course, is that people will insist on coming along and trying to put things in it.",
		"Time is a drug. Too much of it kills you.",
		"In the beginning there was nothing, which exploded.",
		"Five exclamation marks, the sure sign of an insane mind.",
		"The whole of life is just like watching a film.",
		"Fantasy is the impossible made probable.",
		"Real stupidity beats artificial intelligence every time.",
		"A good bookshop is just a genteel Black Hole that knows how to read.",
		"The pen is mightier than the sword if the sword is very short, and the pen is very sharp.",
		"Give a man a fire and he's warm for a day, but set fire to him and he's warm for the rest of his life.",
	}
	
	authors := []string{
		"Terry Pratchett, Diggers",
		"Terry Pratchett, Small Gods", 
		"Terry Pratchett, Lords and Ladies",
		"Terry Pratchett, Reaper Man",
		"Terry Pratchett, Moving Pictures",
		"Terry Pratchett, The Color of Magic",
		"Terry Pratchett, Hogfather",
		"Terry Pratchett, Good Omens",
		"Terry Pratchett, Guards! Guards!",
		"Terry Pratchett, Night Watch",
	}
	
	builder.WriteString("[\n")
	
	for i := 0; i < quoteCount; i++ {
		if i > 0 {
			builder.WriteString(",\n")
		}
		
		quote := quotes[i%len(quotes)]
		author := authors[i%len(authors)]
		
		builder.WriteString("  {\n")
		builder.WriteString(`    "text": "`)
		builder.WriteString(quote)
		builder.WriteString(`",` + "\n")
		builder.WriteString(`    "author": "`)
		builder.WriteString(author)
		builder.WriteString(`"` + "\n")
		builder.WriteString("  }")
	}
	
	builder.WriteString("\n]")
	return builder.String()
}