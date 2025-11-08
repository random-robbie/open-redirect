# Multi-stage build for smaller image size
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o open-redirect .

# Final stage - use chromium image for headless browser
FROM chromedp/headless-shell:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/open-redirect .

# Copy payload file
COPY payloads.txt ./

# Create volume mount point for input/output files
VOLUME ["/app/data"]

# Set default command
ENTRYPOINT ["./open-redirect"]

# Default arguments
CMD ["-urls", "/app/data/urls.txt", "-payloads", "/app/payloads.txt", "-output", "/app/data/found.txt"]
