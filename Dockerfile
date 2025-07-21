# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git=2.45.2-r0 ca-certificates=20240705-r0

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cursor-admin-api-exporter .

# Final stage
FROM alpine:3.20

RUN apk --no-cache add ca-certificates=20240705-r0 wget=1.24.5-r0

# Create app directory with proper permissions for nobody user
RUN mkdir -p /app && \
    chown -R 65534:65534 /app

WORKDIR /app

# Copy the binary from builder and set permissions
COPY --from=builder --chown=65534:65534 /app/cursor-admin-api-exporter .
RUN chmod +x cursor-admin-api-exporter

# Switch to non-root user (nobody)
USER 65534

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Run the binary
ENTRYPOINT ["./cursor-admin-api-exporter"]