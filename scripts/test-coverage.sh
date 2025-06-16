#!/bin/bash

# Personal Disorganiser Test Coverage Script
# Runs tests with coverage reporting and generates HTML report

set -e

echo "🧪 Running tests with coverage..."

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Check if coverage file was generated
if [ ! -f coverage.out ]; then
    echo "❌ Coverage file not generated"
    exit 1
fi

# Generate HTML coverage report
echo "📊 Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html

# Calculate total coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "📈 Total Coverage: ${COVERAGE}%"

# Check coverage threshold
THRESHOLD=80
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo "⚠️  Coverage ${COVERAGE}% is below threshold of ${THRESHOLD}%"
    echo "🎯 Target: Increase test coverage to meet the ${THRESHOLD}% threshold"
    exit 1
else
    echo "✅ Coverage ${COVERAGE}% meets threshold of ${THRESHOLD}%"
fi

# Show per-package coverage breakdown
echo ""
echo "📋 Coverage breakdown by package:"
go tool cover -func=coverage.out | grep -v "total:" | sort -k3 -nr

echo ""
echo "📄 Coverage report saved to: coverage.html"
echo "📄 Coverage data saved to: coverage.out"
echo ""
echo "🌐 Open coverage.html in your browser to view detailed coverage report"