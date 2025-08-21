FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
RUN touch go.sum
COPY go.sum ./

RUN go mod download

COPY . .

# Build binary
RUN go build -o server ./cmd/server

# Final production image
FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/server /app/server

# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*


EXPOSE 8080
CMD ["/app/server"]
