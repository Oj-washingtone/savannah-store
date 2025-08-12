# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for go mod if needed
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go app statically
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .

# Copy any config/migrations if needed
# COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

# Run the binary
CMD ["./server"]