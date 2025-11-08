# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copy go mod files
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copy source code
COPY api/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags="-w -s -X regulation/internal/config.Version=${VERSION:-dev}" \
    -o regulation ./cmd/regulation

# Final stage
FROM scratch

# Copy certificates and timezone data
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /build/regulation /regulation

# Set the entrypoint
ENTRYPOINT ["/regulation"]
