# Use a minimal alpine linux image
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Set environment variables to build for Linux AMD64
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Install dependencies (replace with your actual build command if needed)
RUN go mod download

# Build the application
RUN go build -o main ./cmd

# Use a slimmer image for runtime
FROM alpine:latest

# Copy the binary
COPY --from=builder /app/main /app/main

# Expose port
EXPOSE 443

# Define command to run
CMD ["/app/main"]