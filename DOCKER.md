# WooCommerce MCP Server - Docker Setup

This document describes the Docker setup for the WooCommerce MCP server, following the same pattern used in the chatbot-service.

## 📁 Project Structure

```
woocommerce_mcp/
├── ops/
│   └── docker/
│       ├── dev/
│       │   └── Dockerfile              # Development Dockerfile with Air hot-reload
│       ├── prod/
│       │   └── Dockerfile              # Production multi-stage Dockerfile
│       ├── docker-compose.yml          # Base compose configuration
│       ├── docker-compose.dev.yml      # Development environment
│       ├── docker-compose.prod.yml     # Production environment
│       └── docker-compose.stage.yml    # Staging environment
├── .air.toml                           # Air configuration for hot-reload
├── .dockerignore                       # Docker ignore file
└── Makefile                            # Build and deployment commands
```

## 🚀 Quick Start

### Prerequisites

1. **Traefik Orchestrator Running**: Make sure the `chatbot-dev` orchestrator is running:
   ```bash
   cd ../chatbot-dev
   docker-compose up -d
   ```

2. **Shared Network**: The `shared_network` should be created by the orchestrator.

### Development Environment

```bash
# Start development environment with hot-reload
make docker-dev

# View logs
make docker-logs-dev

# Stop development environment
make docker-dev-down
```

### Production Environment

```bash
# Start production environment
make docker-prod

# View logs
make docker-logs-prod

# Stop production environment
make docker-prod-down
```

### Staging Environment

```bash
# Start staging environment
make docker-stage

# View logs
make docker-logs-stage

# Stop staging environment
make docker-stage-down
```

## 🌐 Access URLs

| Environment | URL | Description |
|-------------|-----|-------------|
| Development | `http://woocommerce-mcp.localhost:8000` | Dev environment with hot-reload |
| Staging | `http://woocommerce-mcp-stage.localhost:8000` | Staging environment |
| Production | `http://woocommerce-mcp.localhost:8000` | Production environment |

## 🔧 Available Endpoints

- `GET /health` - Health check endpoint
- `GET /list_tools` - List available MCP tools
- `POST /call_tool` - Execute MCP tool

### Example Usage

```bash
# Health check
curl http://woocommerce-mcp.localhost:8000/health

# List tools
curl http://woocommerce-mcp.localhost:8000/list_tools

# Search products
curl -X POST http://woocommerce-mcp.localhost:8000/call_tool \
  -H "Content-Type: application/json" \
  -d '{
    "name": "search_products",
    "arguments": {
      "base_url": "https://your-woocommerce-store.com",
      "consumer_key": "ck_your_key",
      "consumer_secret": "cs_your_secret",
      "search": "shirt",
      "per_page": "5"
    }
  }'
```

## 🏗️ Docker Configuration

### Development Dockerfile Features

- **Base Image**: `golang:1.25-alpine`
- **Hot Reload**: Air for automatic rebuilds on code changes
- **Debugging**: Delve debugger support
- **Volume Mounting**: Source code mounted for live development

### Production Dockerfile Features

- **Multi-stage Build**: Optimized for size and security
- **Minimal Runtime**: Alpine Linux base
- **Health Checks**: Built-in health monitoring
- **Security**: Non-root user, minimal attack surface

### Traefik Integration

The service is automatically discovered by Traefik with the following labels:

```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.woocommerce-mcp-dev.rule=Host(`woocommerce-mcp.localhost`)"
  - "traefik.http.services.woocommerce-mcp-dev.loadbalancer.server.port=8080"
  - "traefik.docker.network=shared_network"
```

## 🔍 Monitoring & Debugging

### View Container Logs

```bash
# Development logs
docker logs docker-woocommerce-mcp-1 -f

# Or using make
make docker-logs-dev
```

### Health Checks

The service includes automatic health checks:

```bash
# Check health through Traefik
curl http://woocommerce-mcp.localhost:8000/health

# Direct container health check
docker ps --filter "name=woocommerce-mcp" --format "table {{.Names}}\t{{.Status}}"
```

### Debugging

For development debugging with Delve:

```bash
# The development container includes Delve debugger
# Expose port 2345 for debugging if needed
```

## 🛠️ Makefile Commands

| Command | Description |
|---------|-------------|
| `make help` | Show available commands |
| `make docker-dev` | Start development environment |
| `make docker-prod` | Start production environment |
| `make docker-stage` | Start staging environment |
| `make docker-logs-dev` | View development logs |
| `make docker-logs-prod` | View production logs |
| `make docker-dev-down` | Stop development environment |
| `make docker-prod-down` | Stop production environment |
| `make health` | Check service health |
| `make start-with-orchestrator` | Start with chatbot-dev orchestrator |

## 🔧 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `GIN_MODE` | `debug` (dev) / `release` (prod) | Gin framework mode |

## 🐛 Troubleshooting

### Service Not Accessible Through Traefik

1. **Check if Traefik is running**:
   ```bash
   docker ps | grep traefik
   ```

2. **Verify shared network**:
   ```bash
   docker network ls | grep shared_network
   ```

3. **Check Traefik dashboard**:
   Visit `http://localhost:8080` to see registered services

### Container Won't Start

1. **Check logs**:
   ```bash
   make docker-logs-dev
   ```

2. **Verify Docker build**:
   ```bash
   make docker-build-dev
   ```

### Health Check Failing

1. **Test direct container access**:
   ```bash
   docker exec -it docker-woocommerce-mcp-1 curl http://localhost:8080/health
   ```

2. **Check container status**:
   ```bash
   docker ps --filter "name=woocommerce-mcp"
   ```

## 🔄 Integration with Chatbot Orchestrator

This service is designed to work seamlessly with the `chatbot-dev` orchestrator:

1. **Start orchestrator first**:
   ```bash
   cd ../chatbot-dev
   docker-compose up -d
   ```

2. **Start MCP service**:
   ```bash
   cd ../woocommerce_mcp
   make start-with-orchestrator
   ```

The service will automatically:
- Connect to the `shared_network`
- Register with Traefik for routing
- Provide health checks for monitoring
- Support hot-reload in development mode
