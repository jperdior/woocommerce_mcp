#!/bin/bash

# WooCommerce API Debug Script
# This script helps diagnose WooCommerce API connection issues

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CONTAINER_NAME="woocommerce-mcp-woocommerce-mcp-1"
WORDPRESS_CONTAINER="wordpress-dev-wordpress-1"
WORDPRESS_URL_EXTERNAL="http://wordpress.localhost:8000"
WORDPRESS_URL_INTERNAL="http://wordpress-dev-wordpress-1"

# Your API credentials
CONSUMER_KEY="ck_70d8b0501aff3dd8075f4e34bd09296e809dbad6"
CONSUMER_SECRET="cs_0f64755e1342feeeb66db21b151d4c3c34e1c95a"

echo -e "${BLUE}🔍 WooCommerce API Debug & Test${NC}"
echo "=============================================="

# Test 1: Check if WordPress is accessible
echo -e "${YELLOW}1. Testing WordPress accessibility...${NC}"
if curl -s -f "$WORDPRESS_URL_EXTERNAL" > /dev/null; then
    echo -e "${GREEN}✅ WordPress accessible externally at $WORDPRESS_URL_EXTERNAL${NC}"
else
    echo -e "${RED}❌ WordPress not accessible externally${NC}"
fi

if docker exec "$CONTAINER_NAME" curl -s -f "$WORDPRESS_URL_INTERNAL" > /dev/null; then
    echo -e "${GREEN}✅ WordPress accessible from MCP container${NC}"
else
    echo -e "${RED}❌ WordPress not accessible from MCP container${NC}"
fi

echo

# Test 2: Check WooCommerce API endpoint
echo -e "${YELLOW}2. Testing WooCommerce API endpoint...${NC}"
wc_api_response=$(docker exec "$CONTAINER_NAME" curl -s "$WORDPRESS_URL_INTERNAL/wp-json/wc/v3/")
if echo "$wc_api_response" | grep -q "namespace"; then
    echo -e "${GREEN}✅ WooCommerce API endpoint is available${NC}"
    echo "   Available routes: $(echo "$wc_api_response" | jq -r '.routes | keys | join(", ")' 2>/dev/null || echo "Could not parse routes")"
else
    echo -e "${RED}❌ WooCommerce API endpoint not found${NC}"
    echo "   Response: $wc_api_response"
fi

echo

# Test 3: Test API credentials
echo -e "${YELLOW}3. Testing API credentials...${NC}"
echo "   Consumer Key: ${CONSUMER_KEY:0:20}..."
echo "   Consumer Secret: ${CONSUMER_SECRET:0:20}..."

# Test products endpoint
products_response=$(docker exec "$CONTAINER_NAME" curl -s "$WORDPRESS_URL_INTERNAL/wp-json/wc/v3/products?consumer_key=$CONSUMER_KEY&consumer_secret=$CONSUMER_SECRET")

if echo "$products_response" | grep -q "woocommerce_rest_cannot_view"; then
    echo -e "${RED}❌ Authentication failed - insufficient permissions${NC}"
    echo "   Error: $(echo "$products_response" | jq -r '.message' 2>/dev/null || echo "$products_response")"
    echo
    echo -e "${YELLOW}💡 Possible solutions:${NC}"
    echo "   1. Check API key permissions (should be 'Read' or 'Read/Write')"
    echo "   2. Regenerate the API key"
    echo "   3. Make sure the API key user has admin privileges"
    echo "   4. Check if WooCommerce is properly activated"
elif echo "$products_response" | grep -q "woocommerce_rest_authentication_error"; then
    echo -e "${RED}❌ Authentication error - invalid credentials${NC}"
    echo "   Error: $(echo "$products_response" | jq -r '.message' 2>/dev/null || echo "$products_response")"
elif echo "$products_response" | jq -e '. | type == "array"' > /dev/null 2>&1; then
    echo -e "${GREEN}✅ API credentials are working!${NC}"
    product_count=$(echo "$products_response" | jq '. | length')
    echo "   Found $product_count products in your store"
    
    if [ "$product_count" -gt 0 ]; then
        echo -e "${BLUE}📦 Sample products:${NC}"
        echo "$products_response" | jq -r '.[] | "   - \(.name // "Unnamed") (ID: \(.id), Price: $\(.price // "N/A"))"' | head -5
    else
        echo -e "${YELLOW}⚠️  No products found in your store${NC}"
        echo "   Add some products to test search functionality"
    fi
else
    echo -e "${RED}❌ Unexpected response${NC}"
    echo "   Response: $products_response"
fi

echo

# Test 4: Test MCP server with corrected URL
if echo "$products_response" | jq -e '. | type == "array"' > /dev/null 2>&1; then
    echo -e "${YELLOW}4. Testing MCP server with internal WordPress URL...${NC}"
    
    mcp_request='{
        "name": "search_products",
        "arguments": {
            "base_url": "'$WORDPRESS_URL_INTERNAL'",
            "consumer_key": "'$CONSUMER_KEY'",
            "consumer_secret": "'$CONSUMER_SECRET'",
            "search": "pink",
            "per_page": "5"
        }
    }'
    
    mcp_response=$(docker exec "$CONTAINER_NAME" curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$mcp_request" \
        http://localhost:8080/call_tool)
    
    echo "MCP Response:"
    echo "$mcp_response" | jq '.'
    
    if echo "$mcp_response" | jq -e '.isError == false' > /dev/null 2>&1; then
        echo -e "${GREEN}✅ MCP server is working correctly!${NC}"
        
        # Try to parse and display products
        content=$(echo "$mcp_response" | jq -r '.content[0].text // ""')
        if [ -n "$content" ] && [ "$content" != "null" ]; then
            echo -e "${BLUE}🔍 Search results for 'pink':${NC}"
            if echo "$content" | jq -e '. | type == "array"' > /dev/null 2>&1; then
                product_count=$(echo "$content" | jq '. | length')
                if [ "$product_count" -gt 0 ]; then
                    echo "$content" | jq -r '.[] | "🛍️  \(.name // "Unnamed Product")
   💰 Price: $\(.price // "N/A")
   📋 Status: \(.status // "N/A")
   📦 Stock: \(.stock_status // "N/A")
   🆔 ID: \(.id // "N/A")\n"'
                else
                    echo "   No products found containing 'pink'"
                fi
            else
                echo "   Raw response: $content"
            fi
        fi
    else
        echo -e "${RED}❌ MCP server error${NC}"
        error_msg=$(echo "$mcp_response" | jq -r '.content[0].text // "Unknown error"')
        echo "   Error: $error_msg"
    fi
else
    echo -e "${YELLOW}4. Skipping MCP test due to API credential issues${NC}"
fi

echo
echo -e "${BLUE}🔧 Troubleshooting Guide:${NC}"
echo
echo -e "${YELLOW}If you're getting authentication errors:${NC}"
echo "1. Go to: $WORDPRESS_URL_EXTERNAL/wp-admin/admin.php?page=wc-settings&tab=advanced&section=keys"
echo "2. Delete the existing API key"
echo "3. Create a new one with these settings:"
echo "   - Description: 'MCP Server'"
echo "   - User: Select an Administrator user"
echo "   - Permissions: 'Read' or 'Read/Write'"
echo "4. Copy the new Consumer Key and Consumer Secret"
echo "5. Update your test scripts with the new credentials"
echo
echo -e "${YELLOW}If WooCommerce API is not available:${NC}"
echo "1. Make sure WooCommerce plugin is installed and activated"
echo "2. Go to WooCommerce > Settings > Advanced > REST API"
echo "3. Make sure 'Enable the REST API' is checked"
echo
echo -e "${YELLOW}If no products are found:${NC}"
echo "1. Add some test products in WooCommerce"
echo "2. Make sure products are 'Published' (not Draft)"
echo "3. Try searching for different terms"

echo
echo "Debug completed!"
