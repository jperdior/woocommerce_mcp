package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WooCommerceClient struct {
	BaseURL        string
	ConsumerKey    string
	ConsumerSecret string
	HTTPClient     *http.Client
}

func NewWooCommerceClient(config WooCommerceConfig) *WooCommerceClient {
	return &WooCommerceClient{
		BaseURL:        strings.TrimSuffix(config.BaseURL, "/"),
		ConsumerKey:    config.ConsumerKey,
		ConsumerSecret: config.ConsumerSecret,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (wc *WooCommerceClient) SearchProducts(params map[string]interface{}) ([]Product, error) {
	// Build URL with query parameters
	endpoint := fmt.Sprintf("%s/wp-json/wc/v3/products", wc.BaseURL)

	// Parse base URL
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	// Add query parameters
	query := u.Query()

	// Add authentication parameters
	query.Set("consumer_key", wc.ConsumerKey)
	query.Set("consumer_secret", wc.ConsumerSecret)

	// Add search parameters
	for key, value := range params {
		if value != nil {
			query.Set(key, fmt.Sprintf("%v", value))
		}
	}

	u.RawQuery = query.Encode()

	// Make HTTP request
	resp, err := wc.HTTPClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WooCommerce API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var products []Product
	if err := json.Unmarshal(body, &products); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return products, nil
}
