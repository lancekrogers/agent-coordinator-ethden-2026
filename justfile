@default:
    just --list --justfile {{source_file()}}

build:
    go build -o bin/ ./cmd/...

run *ARGS:
    go run ./cmd/... {{ARGS}}

test *ARGS:
    go test ./... {{ARGS}}

test-cover:
    go test -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

lint:
    golangci-lint run ./...

tidy:
    go mod tidy

clean:
    rm -rf bin/ coverage.out coverage.html
