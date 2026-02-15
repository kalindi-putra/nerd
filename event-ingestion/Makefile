.PHONY: help build up down logs clean restart rebuild test

# Default target
help:
	@echo "Event Ingestion - Docker Deployment Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build Docker images"
	@echo "  up           Start services (detached mode)"
	@echo "  down         Stop and remove containers"
	@echo "  logs         View logs from all services"
	@echo "  logs-server  View server logs"
	@echo "  logs-worker  View worker logs"
	@echo "  restart      Restart all services"
	@echo "  rebuild      Rebuild images and restart services"
	@echo "  clean        Remove containers, networks, and images"
	@echo "  ps           Show running containers"
	@echo "  test         Run a quick test"

# Build Docker images
build:
	@echo "Building Docker images..."
	cd deploy && docker-compose build

# Start services in detached mode
up:
	@echo "Starting services..."
	cd deploy && docker-compose up -d

# Start services in foreground
up-fg:
	@echo "Starting services (foreground)..."
	cd deploy && docker-compose up

# Stop services
down:
	@echo "Stopping services..."
	cd deploy && docker-compose down

# View all logs
logs:
	cd deploy && docker-compose logs -f

# View server logs
logs-server:
	cd deploy && docker-compose logs -f event-server

# View worker logs
logs-worker:
	cd deploy && docker-compose logs -f event-worker

# Restart services
restart:
	@echo "Restarting services..."
	cd deploy && docker-compose restart

# Rebuild and restart
rebuild:
	@echo "Rebuilding and restarting services..."
	cd deploy && docker-compose down
	cd deploy && docker-compose build --no-cache
	cd deploy && docker-compose up -d

# Show running containers
ps:
	cd deploy && docker-compose ps

# Clean up everything
clean:
	@echo "Cleaning up containers, networks, and images..."
	cd deploy && docker-compose down --rmi all -v

# Run a quick test (shows if containers are running)
test:
	@echo "Testing deployment..."
	cd deploy && docker-compose ps
	@echo ""
	@echo "Checking server logs for startup message..."
	cd deploy && docker-compose logs event-server | grep "grpc server running" || echo "Server may still be starting..."
