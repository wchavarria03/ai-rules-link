# Stage 1: Build the Go binary
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# -ldflags="-w -s" strips debug information, making the binary smaller
# CGO_ENABLED=0 is important for creating a static binary that can run in a minimal container
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /ai-rules-link .

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Copy the built binary from the builder stage
COPY --from=builder /ai-rules-link /usr/local/bin/

# The binary is now in the PATH, so it can be run directly
ENTRYPOINT ["ai-rules-link"]
