# WooCommerce MCP - Chatbot Service Integration

This guide explains how the WooCommerce MCP server integrates with your existing chatbot service architecture.

## Architecture Overview

```
[Frontend] → [Traefik] → [Message API] → [MCP Client] → [WooCommerce MCP HTTP Bridge] → [WooCommerce API]
                      ↘ [Chatbot API] ↗
                      ↘ [Database] ↗
                      ↘ [RabbitMQ] ↗
                      ↘ [Redis] ↗
```

## Integration Details

### 1. **Network Integration**
- **Shared Network**: `shared_network` (external) 
- **Port**: `8090` (internal), routed via Traefik
- **URL**: `http://woocommerce-mcp.localhost:8000/` (development)

### 2. **Protocol Compatibility**
The WooCommerce MCP HTTP Bridge implements:
- **JSON-RPC 2.0** over HTTP (compatible with existing chatbot MCP client)
- **Server-Sent Events (SSE)** responses
- **Standard MCP Methods**: `tools/list`, `tools/call`

### 3. **Service Configuration**

#### Development Environment
```yaml
# In chatbot-service/ops/docker/docker-compose.dev.yml
woocommerce-mcp:
  build:
    context: ../../../woocommerce_mcp
    dockerfile: ops/docker/Dockerfile.dev
  environment:
    - PORT=8090
    - ENVIRONMENT=development
  labels:
    - "traefik.enable=true"
    - "traefik.http.routers.woocommerce-mcp.rule=Host(`woocommerce-mcp.localhost`)"
```

#### Production Environment
```yaml
# In chatbot-service/ops/docker/docker-compose.yml
woocommerce-mcp:
  build:
    context: ../../../woocommerce_mcp
    dockerfile: ops/docker/Dockerfile.http-bridge
  environment:
    - PORT=8090
    - ENVIRONMENT=production
  healthcheck:
    test: ["CMD", "curl", "-f", "http://localhost:8090/health"]
```

## Usage in Chatbot Service

### 1. **Adding WooCommerce MCP Server**

In the chatbot frontend, add a new MCP server with:
- **Name**: `WooCommerce MCP`
- **URL**: `http://woocommerce-mcp.localhost:8000/` (development)
- **Description**: `Search and manage WooCommerce products`

### 2. **Available Tools**

#### `search_products`
Search for products in WooCommerce store with advanced filtering.

**Required Parameters:**
- `base_url`: WooCommerce store URL (e.g., `https://mystore.com`)
- `consumer_key`: WooCommerce REST API consumer key
- `consumer_secret`: WooCommerce REST API consumer secret

**Optional Parameters:**
- `search`: Search term for product name/description/SKU
- `category`: Category ID or slug filter
- `tag`: Tag ID or slug filter
- `status`: Product status (`draft`, `pending`, `private`, `publish`)
- `type`: Product type (`simple`, `grouped`, `external`, `variable`)
- `featured`: Featured products filter (`true`/`false`)
- `on_sale`: On sale products filter (`true`/`false`)
- `min_price`: Minimum price filter
- `max_price`: Maximum price filter
- `stock_status`: Stock status (`instock`, `outofstock`, `onbackorder`)
- `per_page`: Items per page (1-100, default: 10)
- `page`: Page number (default: 1)
- `order`: Sort order (`asc`, `desc`)
- `orderby`: Sort field (`date`, `id`, `title`, `slug`, `price`, etc.)

### 3. **Example Requests**

#### Via Chatbot UI
User: "Show me featured sneakers under $100"

The chatbot will automatically:
1. Identify the need for product search
2. Select the `search_products` tool
3. Call WooCommerce MCP with appropriate parameters
4. Return formatted product results

#### Via API (for testing)
```bash
curl -X POST http://woocommerce-mcp.localhost:8000/ \
  -H "Content-Type: application/json" \
  -H "Accept: application/json, text/event-stream" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "id": 1,
    "params": {
      "name": "search_products",
      "arguments": {
        "base_url": "https://mystore.com",
        "consumer_key": "ck_1234567890",
        "consumer_secret": "cs_1234567890",
        "search": "sneakers",
        "featured": "true",
        "max_price": "100",
        "per_page": "5"
      }
    }
  }'
```

## Development Workflow

### Starting Services

**WooCommerce MCP only:**
```bash
# From woocommerce_mcp directory
make start          # Development mode with hot reload
make start-prod     # Production mode
```

**Full stack (manual):**
1. Start base infrastructure: `cd ../chatbot-dev && docker-compose up -d`
2. Start chatbot services: `cd ../chatbot-service && make start`
3. Start WooCommerce MCP: `make start`

### Testing Integration
```bash
# Test health
make health

# Test MCP protocol
make test-integration

# View logs
make logs
```

### Stopping Services
```bash
make stop           # Stop WooCommerce MCP
make stop-prod      # Stop production mode
```

## Configuration

### Environment Variables
- `PORT`: Server port (default: 9000)
- `ENVIRONMENT`: Runtime environment (`development`/`production`)

### WooCommerce Store Setup
1. Generate WooCommerce REST API credentials
2. Configure API permissions (read products)
3. Use credentials in MCP server configuration

### Traefik Routing
- **Development**: `woocommerce-mcp.localhost:8000`
- **Production**: Configure domain in Traefik labels

## Monitoring & Health Checks

### Health Endpoint
```bash
curl http://woocommerce-mcp.localhost:8000/health
```

### Service Status
```bash
docker ps | grep woocommerce-mcp
```

### Logs
```bash
# All services
make chatbot-logs

# WooCommerce MCP only
docker logs chatbot-service_woocommerce-mcp_1 -f
```

## Troubleshooting

### Common Issues

#### 1. Service Not Starting
```bash
# Check Docker logs
docker logs chatbot-service_woocommerce-mcp_1

# Rebuild service
cd ../chatbot-service/ops/docker
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build woocommerce-mcp
```

#### 2. Traefik Routing Issues
```bash
# Check Traefik dashboard
open http://localhost:8080

# Verify labels
docker inspect chatbot-service_woocommerce-mcp_1 | grep traefik
```

#### 3. Network Connectivity
```bash
# Test from another container
docker exec chatbot-service_message-api_1 curl http://woocommerce-mcp:9000/health

# Check shared network
docker network inspect shared_network
```

#### 4. MCP Protocol Issues
```bash
# Test JSON-RPC directly
curl -X POST http://woocommerce-mcp.localhost:8000/ \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'

# Check SSE response
curl -X POST http://woocommerce-mcp.localhost:8000/ \
  -H "Accept: text/event-stream" \
  -d '{"jsonrpc":"2.0","method":"tools/list","id":1}'
```

## Performance Considerations

### Development
- Hot reload enabled for fast development
- Volume mounts for real-time code changes
- Debug logging enabled

### Production
- Multi-stage Docker builds for smaller images
- Health checks for reliability
- Optimized Go builds
- Connection pooling for WooCommerce API

## Security

### API Credentials
- Store WooCommerce credentials securely
- Use environment variables or secret management
- Limit API permissions to read-only where possible

### Network Security
- Services communicate via internal Docker network
- External access only via Traefik
- Health checks don't expose sensitive data

## Next Steps

1. **Configure WooCommerce Store** - Set up REST API credentials
2. **Add MCP Server** - Register in chatbot frontend
3. **Test Integration** - Verify end-to-end functionality
4. **Monitor Performance** - Watch logs and metrics
5. **Scale as Needed** - Add replicas for production load
