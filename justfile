default:
    @just --list

build:
    go build

tidy:
    go mod tidy

run *args:
    go run main.go {{args}}

test dir="...":
    go test ./{{dir}}

coverage:
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out

coverage-html:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

update-dependency *args:
    go get -u ./{{args}}

update-dependencies:
    go get -u ./...

lint:
    golangci-lint run
