# Build stage
FROM golang:1.24-alpine AS builder

# Install git (required for some go modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy all source files
COPY . .

# Download dependencies and generate go.sum
RUN go mod download && go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy .env if exists (for local, but in production use env vars)
COPY --from=builder /app/.env* ./

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]