# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate templ files and build
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -o tact-webui .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary and static files
COPY --from=builder /app/tact-webui .
COPY --from=builder /app/static ./static

EXPOSE 2200

CMD ["./tact-webui"]
