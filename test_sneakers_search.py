#!/usr/bin/env python3
"""
Test script to search for low sneakers using the WooCommerce MCP server.

This script demonstrates how to:
1. Connect to the WooCommerce MCP server
2. List available tools
3. Search for products (sneakers) using the search_products tool
"""

import requests
import json
import sys

# Configuration
MCP_SERVER_URL = "http://localhost:8000"  # Through Traefik
MCP_SERVER_DIRECT_URL = "http://localhost:8080"  # Direct to container (fallback)

# WooCommerce store configuration (you'll need to replace these with your actual values)
WOOCOMMERCE_CONFIG = {
    "base_url": "http://wordpress.localhost:8000",  # Your WordPress/WooCommerce site
    "consumer_key": "your_consumer_key_here",       # Replace with actual key
    "consumer_secret": "your_consumer_secret_here"  # Replace with actual secret
}

def test_mcp_connection(url):
    """Test if the MCP server is accessible."""
    try:
        response = requests.get(f"{url}/health", timeout=5)
        if response.status_code == 200:
            print(f"✅ MCP server is accessible at {url}")
            return True
        else:
            print(f"❌ MCP server returned status {response.status_code} at {url}")
            return False
    except requests.exceptions.RequestException as e:
        print(f"❌ Cannot connect to MCP server at {url}: {e}")
        return False

def list_tools(server_url):
    """List available tools from the MCP server."""
    try:
        response = requests.get(f"{server_url}/list_tools")
        if response.status_code == 200:
            tools = response.json()
            print("📋 Available MCP Tools:")
            for tool in tools.get('tools', []):
                print(f"  - {tool.get('name')}: {tool.get('description')}")
            return tools
        else:
            print(f"❌ Failed to list tools: {response.status_code}")
            return None
    except requests.exceptions.RequestException as e:
        print(f"❌ Error listing tools: {e}")
        return None

def search_sneakers(server_url, woo_config):
    """Search for low sneakers using the WooCommerce MCP server."""
    
    # Prepare the search request
    search_request = {
        "name": "search_products",
        "arguments": {
            # WooCommerce connection details
            "base_url": woo_config["base_url"],
            "consumer_key": woo_config["consumer_key"],
            "consumer_secret": woo_config["consumer_secret"],
            
            # Search parameters for sneakers
            "search": "low sneakers",
            "per_page": "10",
            "status": "publish",
            "type": "simple"
        }
    }
    
    print("🔍 Searching for low sneakers...")
    print(f"Search parameters: {json.dumps(search_request['arguments'], indent=2)}")
    
    try:
        response = requests.post(
            f"{server_url}/call_tool",
            json=search_request,
            headers={"Content-Type": "application/json"},
            timeout=30
        )
        
        if response.status_code == 200:
            result = response.json()
            print("✅ Search completed successfully!")
            
            # Parse and display results
            if result.get('isError'):
                print(f"❌ Search error: {result.get('content', [{}])[0].get('text', 'Unknown error')}")
            else:
                content = result.get('content', [])
                if content and content[0].get('text'):
                    try:
                        products = json.loads(content[0]['text'])
                        print(f"📦 Found {len(products)} products:")
                        
                        for i, product in enumerate(products, 1):
                            print(f"\n{i}. {product.get('name', 'Unnamed Product')}")
                            print(f"   ID: {product.get('id')}")
                            print(f"   Price: ${product.get('price', 'N/A')}")
                            print(f"   Status: {product.get('status', 'N/A')}")
                            print(f"   Stock: {product.get('stock_status', 'N/A')}")
                            if product.get('short_description'):
                                print(f"   Description: {product['short_description'][:100]}...")
                    except json.JSONDecodeError:
                        print("📄 Raw response:")
                        print(content[0]['text'])
                else:
                    print("📄 Empty response or no products found")
            
            return result
            
        else:
            print(f"❌ Search failed with status {response.status_code}")
            print(f"Response: {response.text}")
            return None
            
    except requests.exceptions.RequestException as e:
        print(f"❌ Error during search: {e}")
        return None

def main():
    """Main function to run the sneaker search test."""
    print("🛍️  WooCommerce MCP Sneaker Search Test")
    print("=" * 50)
    
    # Test MCP server connection (try Traefik first, then direct)
    server_url = None
    if test_mcp_connection(MCP_SERVER_URL):
        server_url = MCP_SERVER_URL
    elif test_mcp_connection(MCP_SERVER_DIRECT_URL):
        server_url = MCP_SERVER_DIRECT_URL
        print("⚠️  Using direct container connection (Traefik routing may need fixing)")
    else:
        print("❌ Cannot connect to MCP server. Please ensure it's running.")
        print("\nTo start the MCP server:")
        print("  cd woocommerce_mcp")
        print("  make start")
        sys.exit(1)
    
    print()
    
    # List available tools
    tools = list_tools(server_url)
    if not tools:
        print("❌ Cannot retrieve tools from MCP server")
        sys.exit(1)
    
    print()
    
    # Check if we have the search_products tool
    search_tool = None
    for tool in tools.get('tools', []):
        if tool.get('name') == 'search_products':
            search_tool = tool
            break
    
    if not search_tool:
        print("❌ search_products tool not found in MCP server")
        sys.exit(1)
    
    print("✅ Found search_products tool")
    print()
    
    # Check WooCommerce configuration
    if (WOOCOMMERCE_CONFIG["consumer_key"] == "your_consumer_key_here" or 
        WOOCOMMERCE_CONFIG["consumer_secret"] == "your_consumer_secret_here"):
        print("⚠️  WARNING: Please update the WooCommerce credentials in this script!")
        print("You need to:")
        print("1. Install WooCommerce plugin in your WordPress site")
        print("2. Generate API keys in WooCommerce > Settings > Advanced > REST API")
        print("3. Update WOOCOMMERCE_CONFIG in this script with your keys")
        print()
        print("For testing purposes, we'll try with placeholder credentials...")
        print()
    
    # Perform the sneaker search
    result = search_sneakers(server_url, WOOCOMMERCE_CONFIG)
    
    if result:
        print("\n✅ Test completed successfully!")
    else:
        print("\n❌ Test failed")
        sys.exit(1)

if __name__ == "__main__":
    main()
