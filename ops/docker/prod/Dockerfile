# Build stage
FROM golang:1.25-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Install required dependencies for building
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o woocommerce-mcp .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and curl for health checks
RUN apk --no-cache add ca-certificates curl

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/woocommerce-mcp .

# Expose the port for the MCP server
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the binary
CMD ["./woocommerce-mcp"]
