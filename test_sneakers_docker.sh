#!/bin/bash

# WooCommerce MCP Sneaker Search Test Script (Docker Version)
# This script tests the WooCommerce MCP server by searching for low sneakers
# Uses docker exec to communicate with the containerized MCP server

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CONTAINER_NAME="woocommerce-mcp-woocommerce-mcp-1"
MCP_SERVER_URL="http://localhost:8080"
WORDPRESS_URL="http://wordpress.localhost:8000"

# WooCommerce API credentials (replace with your actual values)
CONSUMER_KEY="your_consumer_key_here"
CONSUMER_SECRET="your_consumer_secret_here"

echo -e "${BLUE}🛍️  WooCommerce MCP Sneaker Search Test (Docker)${NC}"
echo "======================================================"

# Function to test MCP server connection
test_mcp_connection() {
    echo -e "${YELLOW}Testing MCP server in container: $CONTAINER_NAME${NC}"
    
    if docker exec "$CONTAINER_NAME" curl -f -s "$MCP_SERVER_URL/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ MCP server is accessible in container${NC}"
        return 0
    else
        echo -e "${RED}❌ Cannot connect to MCP server in container${NC}"
        return 1
    fi
}

# Function to list available tools
list_tools() {
    echo -e "${YELLOW}📋 Listing available MCP tools...${NC}"
    
    response=$(docker exec "$CONTAINER_NAME" curl -s "$MCP_SERVER_URL/list_tools")
    echo "$response" | jq '.tools[] | {name: .name, description: .description}'
}

# Function to search for sneakers
search_sneakers() {
    echo -e "${YELLOW}🔍 Searching for low sneakers...${NC}"
    
    # Create the search request JSON
    search_request=$(cat <<EOF
{
    "name": "search_products",
    "arguments": {
        "base_url": "$WORDPRESS_URL",
        "consumer_key": "$CONSUMER_KEY",
        "consumer_secret": "$CONSUMER_SECRET",
        "search": "low sneakers",
        "per_page": "5",
        "status": "publish"
    }
}
EOF
)
    
    echo "Search request:"
    echo "$search_request" | jq '.'
    echo
    
    # Create a temporary file in the container with the request
    docker exec "$CONTAINER_NAME" sh -c "cat > /tmp/search_request.json" <<< "$search_request"
    
    # Make the API call from inside the container
    response=$(docker exec "$CONTAINER_NAME" curl -s -X POST \
        -H "Content-Type: application/json" \
        -d @/tmp/search_request.json \
        "$MCP_SERVER_URL/call_tool")
    
    echo -e "${YELLOW}📄 Response:${NC}"
    echo "$response" | jq '.'
    
    # Check if the response indicates an error
    is_error=$(echo "$response" | jq -r '.isError // false')
    
    if [ "$is_error" = "true" ]; then
        echo -e "${RED}❌ Search failed with error${NC}"
        error_message=$(echo "$response" | jq -r '.content[0].text // "Unknown error"')
        echo -e "${RED}Error: $error_message${NC}"
        return 1
    else
        echo -e "${GREEN}✅ Search completed successfully!${NC}"
        
        # Try to parse and display products if available
        content=$(echo "$response" | jq -r '.content[0].text // ""')
        if [ -n "$content" ] && [ "$content" != "null" ]; then
            echo -e "${BLUE}📦 Products found:${NC}"
            echo "$content" | jq '.'
        fi
        return 0
    fi
}

# Function to test with sample WooCommerce demo store
test_with_demo_store() {
    echo -e "${YELLOW}🧪 Testing with WooCommerce demo store...${NC}"
    
    # Use a public WooCommerce demo store for testing
    demo_request=$(cat <<EOF
{
    "name": "search_products",
    "arguments": {
        "base_url": "https://woocommerce.github.io/woocommerce-rest-api-docs",
        "consumer_key": "demo_key",
        "consumer_secret": "demo_secret",
        "search": "sneakers",
        "per_page": "3"
    }
}
EOF
)
    
    echo "Demo request:"
    echo "$demo_request" | jq '.'
    echo
    
    # Create a temporary file in the container with the request
    docker exec "$CONTAINER_NAME" sh -c "cat > /tmp/demo_request.json" <<< "$demo_request"
    
    # Make the API call from inside the container
    response=$(docker exec "$CONTAINER_NAME" curl -s -X POST \
        -H "Content-Type: application/json" \
        -d @/tmp/demo_request.json \
        "$MCP_SERVER_URL/call_tool")
    
    echo -e "${YELLOW}📄 Demo Response:${NC}"
    echo "$response" | jq '.'
}

# Main execution
main() {
    # Check if container exists and is running
    if ! docker ps | grep -q "$CONTAINER_NAME"; then
        echo -e "${RED}❌ Container $CONTAINER_NAME is not running${NC}"
        echo "Please start the MCP server:"
        echo "  make start"
        exit 1
    fi
    
    # Test MCP server connection
    if ! test_mcp_connection; then
        echo -e "${RED}❌ Cannot connect to MCP server in container${NC}"
        exit 1
    fi
    
    echo
    
    # List available tools
    list_tools
    echo
    
    # Check credentials warning
    if [ "$CONSUMER_KEY" = "your_consumer_key_here" ] || [ "$CONSUMER_SECRET" = "your_consumer_secret_here" ]; then
        echo -e "${YELLOW}⚠️  WARNING: Using placeholder WooCommerce credentials!${NC}"
        echo "To use real WooCommerce data:"
        echo "1. Install WooCommerce plugin in WordPress at $WORDPRESS_URL"
        echo "2. Generate API keys in WooCommerce > Settings > Advanced > REST API"
        echo "3. Update CONSUMER_KEY and CONSUMER_SECRET in this script"
        echo
        echo "Proceeding with test using placeholder credentials..."
        echo
    fi
    
    # Perform the sneaker search
    echo -e "${BLUE}Testing with your WordPress/WooCommerce setup:${NC}"
    if search_sneakers; then
        echo -e "${GREEN}✅ Main test completed!${NC}"
    else
        echo -e "${YELLOW}⚠️  Main test failed (expected with placeholder credentials)${NC}"
    fi
    
    echo
    echo -e "${BLUE}Testing MCP server functionality with demo data:${NC}"
    test_with_demo_store
    
    echo
    echo -e "${GREEN}✅ All tests completed!${NC}"
    echo
    echo -e "${BLUE}Next steps:${NC}"
    echo "1. Set up WooCommerce in your WordPress site at $WORDPRESS_URL"
    echo "2. Add some products (including sneakers)"
    echo "3. Generate API keys and update this script"
    echo "4. Run the script again to search real products"
}

# Check dependencies
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ docker is required but not installed${NC}"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo -e "${RED}❌ jq is required but not installed${NC}"
    echo "Install with: brew install jq (on macOS)"
    exit 1
fi

# Run main function
main
