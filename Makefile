GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_MOD_VENDOR=$(GO_CMD) mod vendor
DOCKER_COMPOSE=docker-compose

# Default target
all: vendor build docker

# Ensure vendor directory is populated
vendor:
	@echo "Running go mod vendor..."
	$(GO_MOD_VENDOR)

# Build and run Docker Compose environment
up:
	@echo "Building and running Docker Compose..."
	$(DOCKER_COMPOSE) up --build -d