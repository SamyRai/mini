# Use official Go image
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o mini-mcp ./cmd/mini-mcp

# Final stage
FROM alpine:latest

# Install necessary packages for SSH, networking, and Docker
RUN apk --no-cache add ca-certificates openssh-client curl docker-cli

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/mini-mcp /app/mini-mcp

# Make binary executable
RUN chmod +x /app/mini-mcp

# Run as root for Docker socket access
# USER root

# Expose port (if needed for health checks)
EXPOSE 8080

# Set environment variables
ENV LOG_LEVEL=INFO
ENV VERSION=1.0.0

# Run the application
ENTRYPOINT ["/app/mini-mcp"]
CMD ["--log-level", "INFO"]
