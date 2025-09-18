# WooCommerce MCP Integration Guide

This guide explains how to integrate the WooCommerce MCP server with your message application across Docker containers.

## Understanding MCP Communication

Model Context Protocol (MCP) servers typically communicate via:
- **stdin/stdout** (most common)
- **Command transport** (spawning processes)
- **HTTP** (less common, but useful for microservices)

## Available Implementations

We provide two implementations to suit different architectural needs:

### 1. Pure MCP Server (Recommended for MCP clients)
- **File**: `cmd/mcp/main.go`
- **Transport**: stdin/stdout
- **Use case**: Standard MCP clients, AI applications
- **Docker**: Uses `Dockerfile.mcp`

### 2. HTTP Bridge (Recommended for existing HTTP services)
- **File**: `cmd/http-bridge/main.go`
- **Transport**: HTTP endpoints
- **Use case**: Existing microservices, web applications
- **Docker**: Uses `Dockerfile.http-bridge`

## Docker Container Communication

### Option 1: HTTP Bridge (Easiest)

If your message API currently uses HTTP, use the HTTP bridge:

```yaml
# docker-compose.yml
version: '3.8'
services:
  message-api:
    build: ./chatbot-service
    environment:
      - WOOCOMMERCE_MCP_URL=http://woocommerce-mcp-http:8080
    depends_on:
      - woocommerce-mcp-http

  woocommerce-mcp-http:
    build: 
      context: ./woocommerce_mcp
      dockerfile: ops/docker/Dockerfile.http-bridge
    ports:
      - "8080:8080"
```

Your message API makes HTTP requests:
```bash
curl -X POST http://woocommerce-mcp-http:8080/call_tool \
  -H "Content-Type: application/json" \
  -d '{
    "name": "search_products",
    "arguments": {
      "base_url": "https://mystore.com",
      "consumer_key": "ck_123",
      "consumer_secret": "cs_123",
      "search": "sneakers"
    }
  }'
```

### Option 2: MCP Client Integration (Standard MCP)

For proper MCP integration, your message API uses the MCP Go SDK:

```yaml
# docker-compose.yml
version: '3.8'
services:
  message-api:
    build: ./chatbot-service
    depends_on:
      - woocommerce-mcp
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - WOOCOMMERCE_MCP_CONTAINER=woocommerce-mcp

  woocommerce-mcp:
    build:
      context: ./woocommerce_mcp
      dockerfile: ops/docker/Dockerfile.mcp
    container_name: woocommerce-mcp
```

Your message API code:
```go
// Add to your message API
import "github.com/modelcontextprotocol/go-sdk/mcp"

// Connect to MCP server
client := mcp.NewClient(&mcp.Implementation{Name: "message-api", Version: "v1.0.0"}, nil)
transport := &mcp.CommandTransport{
    Command: exec.Command("docker", "exec", "woocommerce-mcp", "./woocommerce-mcp"),
}
session, err := client.Connect(ctx, transport, nil)

// Call tools
result, err := session.CallTool(ctx, &mcp.CallToolParams{
    Name: "search_products",
    Arguments: map[string]any{
        "base_url": "https://mystore.com",
        "consumer_key": "ck_123",
        "consumer_secret": "cs_123",
        "search": "sneakers",
    },
})
```

## Quick Start

### 1. Start with HTTP Bridge (Immediate compatibility)

```bash
# Build and start HTTP bridge
make local-build-http
make local-run-http

# Test it
curl http://localhost:8080/health
curl http://localhost:8080/list_tools
```

### 2. Start with Docker Compose

```bash
# Start both MCP server and HTTP bridge
make mcp-start

# Check logs
make mcp-logs

# Test HTTP bridge
make mcp-test-http

# Stop services
make mcp-stop
```

### 3. Integrate with your Message API

Choose your approach:

**Option A: Keep using HTTP (easy migration)**
- Update your message API to call `http://woocommerce-mcp-http:8080/call_tool`
- No changes needed to your existing HTTP client code

**Option B: Use MCP SDK (standard approach)**
- Add MCP SDK to your message API: `go get github.com/modelcontextprotocol/go-sdk`
- Implement MCP client as shown in `examples/message-api-integration.go`

## Example Integration

See `examples/message-api-integration.go` for a complete example of how your message API would integrate with the MCP server.

## Benefits of Each Approach

### HTTP Bridge
✅ **Pros:**
- Easy migration from existing HTTP services
- Standard HTTP debugging tools work
- Familiar REST API patterns
- Works with any HTTP client

❌ **Cons:**
- Not "true" MCP protocol
- Additional HTTP overhead
- Limited to HTTP transport features

### Pure MCP
✅ **Pros:**
- Official MCP protocol compliance
- Type-safe Go SDK integration  
- Better performance (no HTTP overhead)
- Future-proof for MCP ecosystem

❌ **Cons:**
- Requires MCP SDK integration
- Different debugging approach
- More complex container communication setup

## Migration Path

1. **Start with HTTP Bridge** - Get immediate compatibility
2. **Test thoroughly** - Ensure all functionality works
3. **Gradually migrate** - Add MCP SDK to your message API
4. **Switch transports** - Move from HTTP to stdio transport
5. **Remove HTTP bridge** - Use pure MCP server

## Troubleshooting

### HTTP Bridge Issues
```bash
# Check if service is running
curl http://localhost:8080/health

# Check logs
docker logs woocommerce-mcp-http

# Test tools endpoint
curl http://localhost:8080/list_tools | jq
```

### MCP Server Issues
```bash
# Test MCP server directly
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | ./woocommerce-mcp-server

# Check container logs
docker logs woocommerce-mcp
```

### Container Communication
```bash
# Check if containers can reach each other
docker exec message-api ping woocommerce-mcp-http

# Check docker network
docker network ls
docker network inspect woocommerce-mcp_mcp-network
```

## Next Steps

1. Choose your integration approach (HTTP bridge for easy start)
2. Update your docker-compose.yml
3. Test the integration
4. Consider migrating to pure MCP for long-term benefits

For more details, see the example files and Dockerfiles provided.
