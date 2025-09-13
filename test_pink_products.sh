#!/bin/bash

# WooCommerce MCP Pink Products Search Test
# This script tests the WooCommerce MCP server by searching for products containing "pink"

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

# WooCommerce API credentials - YOU NEED TO UPDATE THESE!
# Go to: http://wordpress.localhost:8000/wp-admin/admin.php?page=wc-settings&tab=advanced&section=keys
CONSUMER_KEY="ck_70d8b0501aff3dd8075f4e34bd09296e809dbad6"
CONSUMER_SECRET="cs_0f64755e1342feeeb66db21b151d4c3c34e1c95a"

echo -e "${BLUE}🛍️  WooCommerce MCP Pink Products Search Test${NC}"
echo "======================================================="

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

# Function to search for pink products
search_pink_products() {
    echo -e "${YELLOW}🔍 Searching for products containing 'pink'...${NC}"
    
    # Create the search request JSON
    search_request=$(cat <<EOF
{
    "name": "search_products",
    "arguments": {
        "base_url": "$WORDPRESS_URL",
        "consumer_key": "$CONSUMER_KEY",
        "consumer_secret": "$CONSUMER_SECRET",
        "search": "pink",
        "per_page": "10",
        "status": "publish"
    }
}
EOF
)
    
    echo "Search request:"
    echo "$search_request" | jq '.'
    echo
    
    # Create a temporary file in the container with the request
    docker exec "$CONTAINER_NAME" sh -c "cat > /tmp/pink_search_request.json" <<< "$search_request"
    
    # Make the API call from inside the container
    response=$(docker exec "$CONTAINER_NAME" curl -s -X POST \
        -H "Content-Type: application/json" \
        -d @/tmp/pink_search_request.json \
        "$MCP_SERVER_URL/call_tool")
    
    echo -e "${YELLOW}📄 Response:${NC}"
    echo "$response" | jq '.'
    
    # Check if the response indicates an error
    is_error=$(echo "$response" | jq -r '.isError // false')
    
    if [ "$is_error" = "true" ]; then
        echo -e "${RED}❌ Search failed with error${NC}"
        error_message=$(echo "$response" | jq -r '.content[0].text // "Unknown error"')
        echo -e "${RED}Error: $error_message${NC}"
        
        # Check if it's an authentication error
        if [[ "$error_message" == *"consumer_key"* ]] || [[ "$error_message" == *"authentication"* ]]; then
            echo -e "${YELLOW}💡 This looks like an authentication error. Please:${NC}"
            echo "1. Go to: $WORDPRESS_URL/wp-admin/admin.php?page=wc-settings&tab=advanced&section=keys"
            echo "2. Create a new API key with Read permissions"
            echo "3. Update CONSUMER_KEY and CONSUMER_SECRET in this script"
        fi
        
        return 1
    else
        echo -e "${GREEN}✅ Search completed successfully!${NC}"
        
        # Try to parse and display products if available
        content=$(echo "$response" | jq -r '.content[0].text // ""')
        if [ -n "$content" ] && [ "$content" != "null" ]; then
            echo -e "${BLUE}📦 Pink products found:${NC}"
            
            # Try to parse as JSON array of products
            if echo "$content" | jq -e '. | type == "array"' > /dev/null 2>&1; then
                product_count=$(echo "$content" | jq '. | length')
                echo -e "${GREEN}Found $product_count products containing 'pink':${NC}"
                echo
                
                echo "$content" | jq -r '.[] | "🛍️  \(.name // "Unnamed Product")
   💰 Price: $\(.price // "N/A")
   📋 Status: \(.status // "N/A")
   📦 Stock: \(.stock_status // "N/A")
   🆔 ID: \(.id // "N/A")
   📝 Description: \(.short_description // "No description" | .[0:100])...\n"'
            else
                echo -e "${BLUE}Raw response:${NC}"
                echo "$content" | jq '.'
            fi
        else
            echo -e "${YELLOW}📄 No products found or empty response${NC}"
        fi
        return 0
    fi
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
        echo -e "${YELLOW}⚠️  You need to set up WooCommerce API credentials!${NC}"
        echo
        echo -e "${BLUE}To get your API credentials:${NC}"
        echo "1. Go to: $WORDPRESS_URL/wp-admin/"
        echo "2. Navigate to: WooCommerce > Settings > Advanced > REST API"
        echo "3. Click 'Add Key'"
        echo "4. Description: 'MCP Server'"
        echo "5. User: Select your admin user"
        echo "6. Permissions: Read"
        echo "7. Click 'Generate API Key'"
        echo "8. Copy the Consumer Key and Consumer Secret"
        echo "9. Update this script with your keys"
        echo
        echo "Proceeding with test using placeholder credentials (will likely fail)..."
        echo
    fi
    
    # Perform the pink product search
    if search_pink_products; then
        echo -e "${GREEN}✅ Test completed successfully!${NC}"
        echo
        echo -e "${BLUE}🎉 Great! Your WooCommerce MCP server is working!${NC}"
        echo "You can now use it to search for any products in your store."
    else
        echo -e "${YELLOW}⚠️  Test failed - likely due to missing API credentials${NC}"
        echo
        echo -e "${BLUE}Next steps:${NC}"
        echo "1. Set up WooCommerce API credentials (see instructions above)"
        echo "2. Update this script with your real credentials"
        echo "3. Run the test again"
    fi
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
