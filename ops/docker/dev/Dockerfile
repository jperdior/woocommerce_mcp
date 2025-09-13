# Use Go 1.25 with Alpine as the base image
FROM golang:1.25-alpine

# Set the working directory inside the container
WORKDIR /app

# Install required dependencies
RUN apk add --no-cache bash git curl

# Copy the Go application source code to the container
COPY . .

# Install Air for live reloading
RUN go install github.com/air-verse/air@latest

# Install Delve for debugging
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Ensure dependencies are downloaded
RUN go mod download

# Expose the port for the MCP server
EXPOSE 8080

# Start the application with Air for live reloading
CMD ["air", "-c", ".air.toml"]
