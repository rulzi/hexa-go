# Build stage
FROM golang:1.26-alpine AS builder

# Install git and ca-certificates (needed for private repos and TLS)
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -S app && adduser -S app -G app

WORKDIR /home/app

# Copy the binary from builder stage
COPY --from=builder --chown=app:app /app/main .

# Create storage directory
RUN mkdir -p /home/app/storage && chown -R app:app /home/app/storage

USER app

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release

# Run the application
CMD ["./main"]
