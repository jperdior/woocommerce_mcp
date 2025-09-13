#!/bin/bash

# WooCommerce MCP Sneaker Search Test Script
# This script tests the WooCommerce MCP server by searching for low sneakers

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
MCP_SERVER_URL="http://woocommerce-mcp.localhost:8000"
MCP_SERVER_DIRECT_URL="http://localhost:8080"
WORDPRESS_URL="http://wordpress.localhost:8000"

# WooCommerce API credentials (replace with your actual values)
CONSUMER_KEY="your_consumer_key_here"
CONSUMER_SECRET="your_consumer_secret_here"

echo -e "${BLUE}🛍️  WooCommerce MCP Sneaker Search Test${NC}"
echo "=================================================="

# Function to test MCP server connection
test_mcp_connection() {
    local url=$1
    echo -e "${YELLOW}Testing MCP server at: $url${NC}"
    
    if curl -f -s "$url/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ MCP server is accessible at $url${NC}"
        return 0
    else
        echo -e "${RED}❌ Cannot connect to MCP server at $url${NC}"
        return 1
    fi
}

# Function to list available tools
list_tools() {
    local server_url=$1
    echo -e "${YELLOW}📋 Listing available MCP tools...${NC}"
    
    response=$(curl -s "$server_url/list_tools")
    echo "$response" | jq '.tools[] | {name: .name, description: .description}'
}

# Function to search for sneakers
search_sneakers() {
    local server_url=$1
    
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
    
    # Make the API call
    response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$search_request" \
        "$server_url/call_tool")
    
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

# Main execution
main() {
    # Test MCP server connection
    server_url=""
    if test_mcp_connection "$MCP_SERVER_URL"; then
        server_url="$MCP_SERVER_URL"
    elif test_mcp_connection "$MCP_SERVER_DIRECT_URL"; then
        server_url="$MCP_SERVER_DIRECT_URL"
        echo -e "${YELLOW}⚠️  Using direct container connection${NC}"
    else
        echo -e "${RED}❌ Cannot connect to MCP server${NC}"
        echo "Please ensure the MCP server is running:"
        echo "  cd woocommerce_mcp"
        echo "  make start"
        exit 1
    fi
    
    echo
    
    # List available tools
    list_tools "$server_url"
    echo
    
    # Check credentials warning
    if [ "$CONSUMER_KEY" = "your_consumer_key_here" ] || [ "$CONSUMER_SECRET" = "your_consumer_secret_here" ]; then
        echo -e "${YELLOW}⚠️  WARNING: Using placeholder WooCommerce credentials!${NC}"
        echo "To use real WooCommerce data:"
        echo "1. Install WooCommerce plugin in WordPress"
        echo "2. Generate API keys in WooCommerce > Settings > Advanced > REST API"
        echo "3. Update CONSUMER_KEY and CONSUMER_SECRET in this script"
        echo
        echo "Proceeding with test using placeholder credentials..."
        echo
    fi
    
    # Perform the sneaker search
    if search_sneakers "$server_url"; then
        echo -e "${GREEN}✅ Test completed successfully!${NC}"
    else
        echo -e "${RED}❌ Test failed${NC}"
        exit 1
    fi
}

# Check dependencies
if ! command -v curl &> /dev/null; then
    echo -e "${RED}❌ curl is required but not installed${NC}"
    exit 1
fi

if ! command -v jq &> /dev/null; then
    echo -e "${RED}❌ jq is required but not installed${NC}"
    echo "Install with: brew install jq (on macOS)"
    exit 1
fi

# Run main function
main
