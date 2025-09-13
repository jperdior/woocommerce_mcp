# Makefile for woocommerce_mcp
# vim: set ft=make ts=8 noet
# Licence MIT

# Variables
# UNAME		:= $(shell uname -s)
PWD = $(shell pwd)
PROJECT_NAME := woocommerce-mcp
DOCKER_COMPOSE=docker-compose -p ${PROJECT_NAME} -f ${PWD}/ops/docker/docker-compose.yml -f ${PWD}/ops/docker/docker-compose.dev.yml
DOCKER_COMPOSE_PROD=docker-compose -p ${PROJECT_NAME} -f ${PWD}/ops/docker/docker-compose.yml -f ${PWD}/ops/docker/docker-compose.prod.yml
DOCKER_COMPOSE_STAGE=docker-compose -p ${PROJECT_NAME} -f ${PWD}/ops/docker/docker-compose.yml -f ${PWD}/ops/docker/docker-compose.stage.yml
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
local-build: ## Build the application locally
	go build -o woocommerce-mcp ./cmd/api

local-run: ## Run the application locally
	go run ./cmd/api

local-run-main: ## Run the application using main.go (backward compatibility)
	go run .

test: ## Run tests
	go test ./...

clean: ## Clean build artifacts
	rm -f woocommerce-mcp
	rm -rf tmp/

# Docker orchestrator integration
start: build run ### Start the application
	@echo "${GREEN}Starting WooCommerce MCP server (development)...${RESET}"

restart: stop start ### Restart the application

build:
	@${DOCKER_COMPOSE} build --no-cache

run:
	@${DOCKER_COMPOSE} up -d

stop: ### Stop the docker containers
	@${DOCKER_COMPOSE} down --remove-orphans

start-prod: build-prod run-prod ### Start the application in production mode

restart-prod: stop-prod start-prod ### Restart the application in production mode

build-prod:
	@${DOCKER_COMPOSE_PROD} build --no-cache

run-prod:
	@${DOCKER_COMPOSE_PROD} up -d

stop-prod: ### Stop the docker containers in production mode
	@${DOCKER_COMPOSE_PROD} down --remove-orphans

start-stage: build-stage run-stage ### Start the application in stage mode

restart-stage: stop-stage start-stage ### Restart the application in stage mode

build-stage:
	@${DOCKER_COMPOSE_STAGE} build

run-stage:
	@${DOCKER_COMPOSE_STAGE} up -d

stop-stage: ### Stop the docker containers in stage mode
	@${DOCKER_COMPOSE_STAGE} down --remove-orphans

# Docker development commands
docker-dev: ## Start development environment with Docker Compose (standalone)
	cd ops/docker && docker-compose -f docker-compose.dev.yml up --build

docker-dev-down: ## Stop development environment (standalone)
	cd ops/docker && docker-compose -f docker-compose.dev.yml down

docker-logs: ## Show Docker logs
	@${DOCKER_COMPOSE} logs -f

docker-logs-dev: ## Show development Docker logs
	cd ops/docker && docker-compose -f docker-compose.dev.yml logs -f

docker-logs-prod: ## Show production Docker logs
	@${DOCKER_COMPOSE_PROD} logs -f

docker-logs-stage: ## Show staging Docker logs
	@${DOCKER_COMPOSE_STAGE} logs -f

# Health check
health: ## Check if the service is healthy
	curl -f http://localhost:8080/health || curl -f http://woocommerce-mcp.localhost:8000/health

# Manifest check
manifest: ## Show MCP manifest information
	curl -s http://localhost:8080/manifest | jq . || echo "Server not running on port 8080"

manifest-file: ## Show full MCP manifest file
	curl -s http://localhost:8080/manifest.json | jq . || echo "Server not running on port 8080"

# Analysis and linting
analysis: ### Run static analysis and linter
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:latest golangci-lint run --fix

# Legacy aliases for backward compatibility
start-with-orchestrator: start ## Alias for start (backward compatibility)

stop-with-orchestrator: stop ## Alias for stop (backward compatibility)
