default:
    @just --list

build:
    go build

tidy:
    go mod tidy

run *args:
    go run main.go {{args}}

# Run tests in a directory (default: all packages)
test dir="...":
    go test -v -timeout=30s ./{{dir}}

# Run tests with verbose output
test-verbose dir="...":
    go test -v -race -timeout=30s -count=1 ./{{dir}}

# Run tests with race detection only (faster)
test-race dir="...":
    go test -race -timeout=30s ./{{dir}}

# Generate coverage report and show summary
coverage:
    @echo "Running tests with coverage..."
    go test -v -race -timeout=30s -coverprofile=coverage.out -covermode=atomic ./...
    @echo ""
    @echo "📊 Coverage Summary:"
    go tool cover -func=coverage.out | tail -1
    @echo ""
    @echo "📋 Detailed Coverage by Function:"
    go tool cover -func=coverage.out

# Generate HTML coverage report
coverage-html:
    @echo "Running tests with coverage..."
    go test -v -race -timeout=30s -coverprofile=coverage.out -covermode=atomic ./...
    @echo ""
    @echo "Generating HTML coverage report..."
    go tool cover -html=coverage.out -o coverage.html
    @echo "✅ Coverage report generated: coverage.html"
    @echo "Open coverage.html in your browser to view detailed coverage"

# Show coverage summary only
coverage-summary:
    @echo "Running tests with coverage..."
    go test -v -race -timeout=30s -coverprofile=coverage.out -covermode=atomic ./...
    @echo ""
    @echo "📊 Coverage Summary:"
    go tool cover -func=coverage.out | tail -1

# Clean coverage files
coverage-clean:
    rm -f coverage.out coverage.html
    @echo "✅ Cleaned coverage files"

update-dependency *args:
    go get -u ./{{args}}

update-dependencies:
    go get -u ./...

lint:
    golangci-lint run
