# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o main ./cmd/app

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy migrations folder (needed for goose)
COPY --from=builder /app/migrations ./migrations

# Expose HTTP and gRPC ports
EXPOSE 8080 50051

CMD ["./main"]