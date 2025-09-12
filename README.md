# WooCommerce MCP Server

A simple Model Context Protocol (MCP) server in Go using the Gin framework for WooCommerce product search functionality.

## Features

- **Search Products**: Search for products in WooCommerce stores using the REST API
- **MCP Protocol**: Implements MCP server endpoints (`/list_tools` and `/call_tool`)
- **Stateless Design**: All configuration (base URL, credentials) provided per request - no server-side storage
- **Authentication**: Supports WooCommerce REST API authentication via consumer key/secret
- **Filtering**: Comprehensive product filtering options (search terms, categories, price, stock status, etc.)
- **Multi-Store Support**: Can work with different WooCommerce stores in the same session

## Installation

1. Clone or navigate to the project directory:
```bash
cd woocommerce_mcp
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the server:
```bash
go build -o woocommerce-mcp
```

## Usage

### Starting the Server

```bash
./woocommerce-mcp
```

The server will start on port 8080 by default. You can set a custom port using the `PORT` environment variable:

```bash
PORT=3000 ./woocommerce-mcp
```

### Available Endpoints

- `GET /health` - Health check endpoint
- `GET /list_tools` - Lists available MCP tools
- `POST /call_tool` - Executes a specific tool

### Search Products Tool

The `search_products` tool allows you to search for products in a WooCommerce store.

#### Required Parameters (provided with each request)

- `base_url`: WooCommerce store base URL (e.g., `https://example.com`)
- `consumer_key`: WooCommerce REST API consumer key  
- `consumer_secret`: WooCommerce REST API consumer secret

> **Note**: The MCP server is stateless - all configuration including the store URL and API credentials must be provided with each request. This allows the same server instance to work with multiple different WooCommerce stores.

#### Optional Parameters

- `search`: Search term to filter products by name, description, or SKU
- `category`: Category ID or slug to filter products
- `tag`: Tag ID or slug to filter products
- `status`: Product status filter (`draft`, `pending`, `private`, `publish`)
- `type`: Product type filter (`simple`, `grouped`, `external`, `variable`)
- `featured`: Filter by featured products (`true`/`false`)
- `on_sale`: Filter by products on sale (`true`/`false`)
- `min_price`: Minimum price filter
- `max_price`: Maximum price filter
- `stock_status`: Stock status filter (`instock`, `outofstock`, `onbackorder`)
- `per_page`: Number of products per page (default: 10, max: 100)
- `page`: Page number for pagination (default: 1)
- `order`: Sort order (`asc`, `desc`)
- `orderby`: Sort by field (`date`, `id`, `include`, `title`, `slug`, `price`, `popularity`, `rating`, `menu_order`)

### Example Usage

#### List Available Tools

```bash
curl http://localhost:8080/list_tools
```

#### Search Products

```bash
curl -X POST http://localhost:8080/call_tool \
  -H "Content-Type: application/json" \
  -d '{
    "name": "search_products",
    "arguments": {
      "base_url": "https://your-woocommerce-store.com",
      "consumer_key": "ck_your_consumer_key",
      "consumer_secret": "cs_your_consumer_secret",
      "search": "shirt",
      "per_page": "5",
      "status": "publish"
    }
  }'
```

#### Search with Category Filter

```bash
curl -X POST http://localhost:8080/call_tool \
  -H "Content-Type: application/json" \
  -d '{
    "name": "search_products",
    "arguments": {
      "base_url": "https://your-woocommerce-store.com",
      "consumer_key": "ck_your_consumer_key",
      "consumer_secret": "cs_your_consumer_secret",
      "category": "clothing",
      "featured": "true",
      "per_page": "10"
    }
  }'
```

## WooCommerce REST API Setup

To use this MCP server, you need to set up REST API access in your WooCommerce store:

1. Go to **WooCommerce > Settings > Advanced > REST API**
2. Click **Add Key**
3. Choose a user and set permissions (usually "Read" is sufficient for product search)
4. Copy the generated Consumer Key and Consumer Secret
5. Use these credentials in your API calls

## Development

### Project Structure

```
woocommerce_mcp/
├── main.go              # Main server entry point
├── types.go             # Type definitions for MCP and WooCommerce
├── handlers.go          # HTTP handlers for MCP endpoints
├── woocommerce_client.go # WooCommerce API client
├── go.mod               # Go module file
└── README.md            # This file
```

### Building

```bash
go build -o woocommerce-mcp
```

### Running in Development

```bash
go run .
```

## Error Handling

The server includes comprehensive error handling for:

- Invalid request formats
- Missing required parameters
- WooCommerce API errors
- Network connectivity issues
- JSON parsing errors

All errors are returned in the MCP format with appropriate error flags.

## Security Considerations

- API credentials are passed with each request and not stored server-side
- CORS is enabled for cross-origin requests
- Use HTTPS in production environments
- Validate and sanitize all input parameters
- Consider rate limiting for production deployments

## License

This project is provided as-is for educational and development purposes.
