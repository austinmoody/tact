.PHONY: build dev generate clean install-templ

# Install templ CLI if not present
install-templ:
	@which templ > /dev/null || go install github.com/a-h/templ/cmd/templ@latest

# Generate templ files
generate: install-templ
	templ generate

# Build the application
build: generate
	go build -o bin/tact-webui .

# Run in development mode with auto-reload
dev: generate
	go run .

# Watch for changes and regenerate (requires templ watch)
watch: install-templ
	templ generate --watch &
	go run .

# Clean build artifacts
clean:
	rm -rf bin/
	find . -name "*_templ.go" -delete

# Run with custom API URL
dev-api: generate
	TACT_API_URL=http://localhost:2100 go run .
