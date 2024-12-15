# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install required system dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/bin/server .
COPY --from=builder /app/config/config.yaml ./config/
COPY --from=builder /app/database/migrations ./database/migrations/

# Expose the application port
EXPOSE 3000

# Run the application
CMD ["./server", "serve", "--config", "config/config.yaml"]
