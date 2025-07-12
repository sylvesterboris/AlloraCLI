# AlloraCLI Dockerfile

# Build stage
FROM golang:1.23-alpine AS builder

# Install git and ca-certificates
RUN apk --no-cache add git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Verify again after copying source
RUN go mod verify

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o allora ./cmd/allora

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S allora && \
    adduser -u 1001 -S allora -G allora

# Set working directory
WORKDIR /home/allora

# Copy the binary from builder stage
COPY --from=builder /app/allora .

# Create config directory
RUN mkdir -p .config/alloracli

# Change ownership
RUN chown -R allora:allora /home/allora

# Switch to non-root user
USER allora

# Expose port (if needed for future web interface)
EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["./allora"]

# Default command
CMD ["--help"]
