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
local-build: ## Build the application locally (HTTP bridge)
	go build -o woocommerce-mcp ./cmd/api

local-build-mcp: ## Build the MCP server locally
	go build -o woocommerce-mcp-server ./cmd/mcp

local-build-http: ## Build the HTTP bridge locally
	go build -o woocommerce-mcp-http ./cmd/http-bridge

local-run: ## Run the application locally (HTTP bridge)
	go run ./cmd/api

local-run-mcp: ## Run the MCP server locally
	go run ./cmd/mcp

local-run-http: ## Run the HTTP bridge locally
	go run ./cmd/http-bridge

local-run-main: ## Run the application using main.go (backward compatibility)
	go run .

test: ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -f woocommerce-mcp woocommerce-mcp-server woocommerce-mcp-http
	rm -rf tmp/

# MCP-specific commands
mcp-start: ## Start MCP services (both server and HTTP bridge)
	docker-compose -p ${PROJECT_NAME} -f ${PWD}/ops/docker/docker-compose.mcp.yml up --build -d

mcp-stop: ## Stop MCP services
	docker-compose -p ${PROJECT_NAME} -f ${PWD}/ops/docker/docker-compose.mcp.yml down --remove-orphans

mcp-logs: ## Show MCP services logs
	docker-compose -p ${PROJECT_NAME} -f ${PWD}/ops/docker/docker-compose.mcp.yml logs -f

mcp-test: ## Test MCP server connection
	@echo "Testing MCP server (stdio)..."
	echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | ./woocommerce-mcp-server || echo "Build MCP server first with 'make local-build-mcp'"

mcp-test-http: ## Test HTTP bridge
	@echo "Testing HTTP bridge..."
	curl -f http://localhost:8090/health && echo " ✓ Health check passed" || echo " ✗ Health check failed"

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
