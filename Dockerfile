# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Install required system dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o /server

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /server /usr/local/bin/server

# Copy config files
COPY config /app/config

EXPOSE 3000

CMD ["/usr/local/bin/server", "serve", "--config", "config/config.yaml"]

