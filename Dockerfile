# ---- Build stage ----
FROM golang:1.23.2-alpine AS builder

# Create working directory
WORKDIR /app

# Copy go.mod and go.sum first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server ./cmd/server

# ---- Runtime stage ----
FROM alpine:3.19

# Create non-root user (security best practice)
RUN adduser -D appuser
USER appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Expose HTTP port
EXPOSE 8080

# Run the service
ENTRYPOINT ["./server"]
