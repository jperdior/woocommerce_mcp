# WooCommerce MCP Server Testing Guide

This guide shows you how to test the WooCommerce MCP server by searching for products like "low sneakers".

## 🚀 Quick Start

### 1. Start the Services

Make sure both WordPress and the WooCommerce MCP server are running:

```bash
# Start the orchestrator (includes WordPress and MCP server)
cd ../chatbot-dev
make start

# Or start individually
cd ../wordpress-dev
make start

cd ../woocommerce_mcp
make start
```

### 2. Set Up WordPress & WooCommerce

1. **Complete WordPress Installation**:
   - Visit: http://wordpress.localhost:8000
   - Follow the installation wizard
   - Create your admin account

2. **Install WooCommerce Plugin**:
   - Log into WordPress admin: http://wordpress.localhost:8000/wp-admin
   - Go to Plugins > Add New
   - Search for "WooCommerce"
   - Install and activate the plugin
   - Follow the WooCommerce setup wizard

3. **Add Some Test Products**:
   - Go to Products > Add New
   - Add a few products including some sneakers:
     - "Low Top Sneakers"
     - "High Top Sneakers" 
     - "Running Shoes"
   - Set prices, descriptions, and publish them

4. **Generate API Keys**:
   - Go to WooCommerce > Settings > Advanced > REST API
   - Click "Add Key"
   - Description: "MCP Server"
   - User: Select your admin user
   - Permissions: Read/Write
   - Click "Generate API Key"
   - **Copy the Consumer Key and Consumer Secret** (you'll need these!)

### 3. Test the MCP Server

#### Option A: Use the Docker Test Script (Recommended)

```bash
# Edit the script to add your API credentials
nano test_sneakers_docker.sh

# Update these lines with your actual keys:
CONSUMER_KEY="ck_your_actual_consumer_key_here"
CONSUMER_SECRET="cs_your_actual_consumer_secret_here"

# Run the test
./test_sneakers_docker.sh
```

#### Option B: Use the Python Test Script

```bash
# Install Python dependencies if needed
pip3 install requests

# Edit the script to add your API credentials
nano test_sneakers_search.py

# Update the WOOCOMMERCE_CONFIG section with your keys

# Run the test
python3 test_sneakers_search.py
```

#### Option C: Manual cURL Test

```bash
# Test the MCP server directly
docker exec woocommerce-mcp-woocommerce-mcp-1 curl -s http://localhost:8080/list_tools | jq '.'

# Search for products (replace with your actual keys)
docker exec woocommerce-mcp-woocommerce-mcp-1 curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "search_products",
    "arguments": {
      "base_url": "http://wordpress.localhost:8000",
      "consumer_key": "ck_your_key_here",
      "consumer_secret": "cs_your_secret_here",
      "search": "low sneakers",
      "per_page": "5"
    }
  }' \
  http://localhost:8080/call_tool | jq '.'
```

## 📋 Available Test Scripts

| Script | Description | Usage |
|--------|-------------|-------|
| `test_sneakers_docker.sh` | **Recommended** - Uses docker exec to communicate with MCP server | `./test_sneakers_docker.sh` |
| `test_sneakers_search.py` | Python version with detailed output | `python3 test_sneakers_search.py` |
| `test_sneakers_search.sh` | Bash version (requires Traefik routing) | `./test_sneakers_search.sh` |

## 🔧 Troubleshooting

### MCP Server Not Accessible

If you get connection errors:

```bash
# Check if containers are running
docker ps | grep -E "(wordpress|woocommerce)"

# Check MCP server health
docker exec woocommerce-mcp-woocommerce-mcp-1 curl http://localhost:8080/health

# Restart MCP server
make restart
```

### WordPress Not Accessible

If WordPress isn't accessible at http://wordpress.localhost:8000:

```bash
# Check Traefik routing
curl -s http://localhost:8080/api/http/routers | jq '.[] | select(.name | contains("wordpress"))'

# Restart WordPress
cd ../wordpress-dev
make restart
```

### WooCommerce API Errors

Common API errors and solutions:

- **"Invalid consumer key"**: Double-check your consumer key and secret
- **"Insufficient permissions"**: Make sure API key has Read/Write permissions
- **"Product not found"**: Add some products to your WooCommerce store first
- **"Invalid request format: EOF"**: Usually means WooCommerce isn't set up or API keys are wrong

## 🛍️ Example Search Parameters

The MCP server supports various search parameters:

```json
{
  "name": "search_products",
  "arguments": {
    "base_url": "http://wordpress.localhost:8000",
    "consumer_key": "your_key",
    "consumer_secret": "your_secret",
    
    // Search parameters
    "search": "sneakers",           // Search term
    "category": "shoes",            // Category slug
    "tag": "athletic",              // Tag slug
    "status": "publish",            // Product status
    "type": "simple",               // Product type
    "featured": "true",             // Featured products only
    "on_sale": "true",              // On sale products only
    "min_price": "50",              // Minimum price
    "max_price": "200",             // Maximum price
    "stock_status": "instock",      // Stock status
    "per_page": "10",               // Results per page (max 100)
    "page": "1",                    // Page number
    "order": "desc",                // Sort order
    "orderby": "date"               // Sort by field
  }
}
```

## 🎯 Expected Results

When working correctly, you should see output like:

```json
{
  "content": [
    {
      "type": "text",
      "text": "[{\"id\":123,\"name\":\"Low Top Sneakers\",\"price\":\"89.99\",\"status\":\"publish\",\"stock_status\":\"instock\"}]"
    }
  ],
  "isError": false
}
```

## 🔗 Useful URLs

- **WordPress**: http://wordpress.localhost:8000
- **WordPress Admin**: http://wordpress.localhost:8000/wp-admin
- **Traefik Dashboard**: http://localhost:8080
- **MCP Server Health**: `docker exec woocommerce-mcp-woocommerce-mcp-1 curl http://localhost:8080/health`

## 📚 Next Steps

Once you have the basic search working:

1. **Add More Products**: Create a variety of products to test different search parameters
2. **Test Categories**: Create product categories and test category-based searches
3. **Test Filters**: Try price ranges, stock status, and other filters
4. **Integration**: Use the MCP server in your applications or AI tools
5. **Extend Functionality**: Add more MCP tools for creating/updating products

## 🆘 Need Help?

If you encounter issues:

1. Check the container logs: `docker logs woocommerce-mcp-woocommerce-mcp-1`
2. Verify WordPress is accessible: `curl http://wordpress.localhost:8000`
3. Test WooCommerce API directly from WordPress admin
4. Make sure all services are running: `docker ps`
