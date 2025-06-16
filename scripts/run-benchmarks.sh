#!/bin/bash

# Personal Disorganiser Benchmark Script
# Runs performance benchmarks for all modules

set -e

echo "🚀 Running Performance Benchmarks for Personal Disorganiser"
echo "============================================================"

# Create benchmarks directory if it doesn't exist
mkdir -p benchmarks

# Function to run benchmarks for a module
run_module_benchmarks() {
    local module=$1
    local module_name=$2
    
    echo ""
    echo "📊 Running benchmarks for $module_name module..."
    echo "------------------------------------------------"
    
    # Run benchmarks and save results
    go test -bench=. -benchmem -run=^$ ./$module > benchmarks/${module_name}_bench.txt 2>&1 || {
        echo "⚠️  No benchmarks found for $module_name module"
        echo "No benchmarks available" > benchmarks/${module_name}_bench.txt
    }
    
    # Display results
    if grep -q "BenchmarkEngine\|BenchmarkStorage\|BenchmarkParse\|BenchmarkLoad" benchmarks/${module_name}_bench.txt; then
        echo "✅ Benchmarks completed for $module_name"
        grep "Benchmark" benchmarks/${module_name}_bench.txt | head -10
    else
        echo "ℹ️  No benchmarks available for $module_name"
    fi
}

# Run benchmarks for each module
run_module_benchmarks "internal/storage" "storage"
run_module_benchmarks "internal/parser" "parser" 
run_module_benchmarks "internal/search" "search"
run_module_benchmarks "internal/quotes" "quotes"
run_module_benchmarks "internal/calendar" "calendar"
run_module_benchmarks "internal/theme" "theme"
run_module_benchmarks "internal/app" "app"

# Generate summary report
echo ""
echo "📋 Generating Benchmark Summary Report..."
echo "=========================================="

cat > benchmarks/summary.md << 'EOF'
# Personal Disorganiser - Performance Benchmark Results

This report contains performance benchmark results for all modules in the Personal Disorganiser application.

## Benchmark Overview

The benchmarks test the performance of core operations across different data sizes:
- **Small**: 10 items
- **Medium**: 100 items  
- **Large**: 1,000 items
- **Extra Large**: 10,000 items (where applicable)

## Results by Module

EOF

# Add results for each module
for module in storage parser search quotes calendar theme app; do
    echo "### ${module^} Module" >> benchmarks/summary.md
    echo "" >> benchmarks/summary.md
    
    if [ -f "benchmarks/${module}_bench.txt" ] && grep -q "Benchmark" benchmarks/${module}_bench.txt; then
        echo '```' >> benchmarks/summary.md
        grep "Benchmark" benchmarks/${module}_bench.txt >> benchmarks/summary.md
        echo '```' >> benchmarks/summary.md
    else
        echo "No benchmarks available for this module." >> benchmarks/summary.md
    fi
    echo "" >> benchmarks/summary.md
done

# Add interpretation guide
cat >> benchmarks/summary.md << 'EOF'
## Interpreting Results

- **ns/op**: Nanoseconds per operation (lower is better)
- **B/op**: Bytes allocated per operation (lower is better)  
- **allocs/op**: Number of allocations per operation (lower is better)

## Performance Guidelines

- Operations under 1µs (1000 ns) are considered very fast
- Operations under 10µs (10000 ns) are considered fast
- Operations under 100µs (100000 ns) are considered acceptable
- Operations over 1ms (1000000 ns) may need optimization

## Generated On

EOF

date >> benchmarks/summary.md

echo ""
echo "✅ Benchmark summary generated: benchmarks/summary.md"
echo ""
echo "🔍 Performance Analysis Complete!"
echo "=================================="
echo "• Individual results: benchmarks/*_bench.txt"
echo "• Summary report: benchmarks/summary.md"
echo "• Coverage report: coverage.html"
echo ""
echo "📈 Review the results to identify optimization opportunities"