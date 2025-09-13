# WooCommerce API Setup for MCP Server

## 🚀 Quick Setup Guide

Your WooCommerce MCP server is working perfectly! You just need to connect it to your WooCommerce store.

### Step 1: Generate WooCommerce API Keys

1. **Go to your WordPress admin**: http://wordpress.localhost:8000/wp-admin/

2. **Navigate to WooCommerce API settings**:
   - Go to: **WooCommerce** → **Settings** → **Advanced** → **REST API**
   - Or directly: http://wordpress.localhost:8000/wp-admin/admin.php?page=wc-settings&tab=advanced&section=keys

3. **Create a new API key**:
   - Click **"Add Key"**
   - **Description**: `MCP Server`
   - **User**: Select your admin user
   - **Permissions**: `Read` (or `Read/Write` if you want to create/update products later)
   - Click **"Generate API Key"**

4. **Copy your credentials**:
   - **Consumer Key**: `ck_xxxxxxxxxxxxxxxxxx`
   - **Consumer Secret**: `cs_xxxxxxxxxxxxxxxxxx`
   - ⚠️ **Important**: Copy these now - you won't see the secret again!

### Step 2: Update the Test Script

1. **Edit the test script**:
   ```bash
   nano test_pink_products.sh
   ```

2. **Update these lines** (around line 15-16):
   ```bash
   # Replace these with your actual keys:
   CONSUMER_KEY="ck_your_actual_consumer_key_here"
   CONSUMER_SECRET="cs_your_actual_consumer_secret_here"
   ```

3. **Save the file** (Ctrl+X, then Y, then Enter)

### Step 3: Test the Search

1. **Run the test**:
   ```bash
   ./test_pink_products.sh
   ```

2. **Expected result**: The script will search for products containing "pink" in your WooCommerce store.

## 🛍️ Add Some Test Products (Optional)

If you don't have products yet, add some test products:

1. Go to: **Products** → **Add New**
2. Create products like:
   - "Pink T-Shirt" - $25.00
   - "Pink Sneakers" - $89.99  
   - "Rose Pink Dress" - $45.00
   - "Hot Pink Phone Case" - $15.99

## 🔧 Manual Test (Alternative)

If you prefer to test manually with curl:

```bash
# Replace YOUR_KEY and YOUR_SECRET with your actual credentials
docker exec woocommerce-mcp-woocommerce-mcp-1 curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "name": "search_products",
    "arguments": {
      "base_url": "http://wordpress.localhost:8000",
      "consumer_key": "YOUR_KEY",
      "consumer_secret": "YOUR_SECRET",
      "search": "pink",
      "per_page": "5"
    }
  }' \
  http://localhost:8080/call_tool | jq '.'
```

## 🎯 What You'll See When It Works

When successful, you'll see output like:

```json
{
  "content": [
    {
      "type": "text", 
      "text": "[{\"id\":123,\"name\":\"Pink T-Shirt\",\"price\":\"25.00\",\"status\":\"publish\"}]"
    }
  ],
  "isError": false
}
```

## 🔍 Available Search Parameters

Your MCP server supports these search options:

- **search**: "pink", "shoes", "dress", etc.
- **category**: Category slug or ID
- **tag**: Tag slug or ID  
- **min_price** / **max_price**: Price range
- **status**: "publish", "draft", etc.
- **stock_status**: "instock", "outofstock"
- **featured**: "true" for featured products only
- **on_sale**: "true" for sale products only
- **per_page**: Number of results (max 100)
- **orderby**: "date", "price", "title", etc.

## 🚨 Troubleshooting

### "Invalid consumer key" error
- Double-check your consumer key and secret
- Make sure you copied them correctly (no extra spaces)

### "Insufficient permissions" error  
- Make sure your API key has "Read" permissions
- Try regenerating the API key

### "No products found"
- Add some products to your WooCommerce store first
- Make sure products are published (not draft)

### Connection errors
- Make sure both WordPress and MCP server are running:
  ```bash
  docker ps | grep -E "(wordpress|woocommerce)"
  ```

## 🎉 Success!

Once you see products in the search results, your WooCommerce MCP server is fully working! You can now:

1. **Use it in AI applications** to search your product catalog
2. **Integrate with chatbots** for customer service
3. **Build custom tools** that interact with your WooCommerce store
4. **Extend functionality** by adding more MCP tools (create products, update inventory, etc.)

---

**Need help?** Check the logs: `docker logs woocommerce-mcp-woocommerce-mcp-1`
