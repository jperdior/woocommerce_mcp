# Makefile for woocommerce_mcp
# vim: set ft=make ts=8 noet
# Licence MIT

# Variables
# UNAME		:= $(shell uname -s)
PWD = $(shell pwd)
PROJECT_NAME := woocommerce-mcp
DOCKER_COMPOSE=docker-compose -p ${PROJECT_NAME} -f ${PWD}/ops/docker/docker-compose.yml -f ${PWD}/ops/docker/docker-compose.dev.yml
DOCKER_COMPOSE_PROD=docker-compose -p ${PROJECT_NAME} -f ${PWD}/ops/docker/docker-compose.yml
GREEN=\033[0;32m
RESET=\033[0m

.EXPORT_ALL_VARIABLES:

# this is godly
# https://news.ycombinator.com/item?id=11939200
.PHONY: help
help:
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@# this is not tested, but prepared in advance for you, Mac drivers
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

# Local development
local-build: ## Build the HTTP bridge locally
	go build -o woocommerce-mcp-http ./cmd/http-bridge

local-run: ## Run the HTTP bridge locally
	go run ./cmd/http-bridge

test: ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -f woocommerce-mcp-http
	rm -rf tmp/

# Docker orchestrator integration (following chatbot-service pattern)
start: build run ## Start WooCommerce MCP (development)
	@echo "${GREEN}Starting WooCommerce MCP server (development)...${RESET}"

restart: stop start ## Restart WooCommerce MCP

build:
	@${DOCKER_COMPOSE} build

run:
	@${DOCKER_COMPOSE} up -d

stop: ## Stop WooCommerce MCP containers
	@${DOCKER_COMPOSE} down --remove-orphans

start-prod: build-prod run-prod ## Start WooCommerce MCP (production)

restart-prod: stop-prod start-prod ## Restart WooCommerce MCP (production)

build-prod:
	@${DOCKER_COMPOSE_PROD} build

run-prod:
	@${DOCKER_COMPOSE_PROD} up -d

stop-prod: ## Stop WooCommerce MCP (production)
	@${DOCKER_COMPOSE_PROD} down --remove-orphans

logs: ## Show WooCommerce MCP logs
	@${DOCKER_COMPOSE} logs -f

logs-prod: ## Show WooCommerce MCP production logs
	@${DOCKER_COMPOSE_PROD} logs -f

test-integration: ## Test WooCommerce MCP integration
	@echo "Testing WooCommerce MCP JSON-RPC endpoint..."
	@curl -X POST http://local.woocommerce-mcp.com:8000/ \
		-H "Content-Type: application/json" \
		-H "Accept: application/json, text/event-stream" \
		-d '{"jsonrpc":"2.0","method":"tools/list","id":1}' \
		|| echo "Make sure WooCommerce MCP is running with 'make start'"

# Health check
health: ## Check if the service is healthy
	curl -f http://localhost:8090/health || curl -f http://local.woocommerce-mcp.com:8000/health

# Manifest check
manifest: ## Show MCP manifest information
	curl -s http://localhost:8090/manifest | jq . || echo "Server not running on port 8090"

manifest-file: ## Show full MCP manifest file
	curl -s http://localhost:8090/manifest.json | jq . || echo "Server not running on port 8090"

# Analysis and linting
analysis: ## Run static analysis and linter
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run --fix
