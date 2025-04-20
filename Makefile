.PHONY: build up down restart logs

# Build Docker containers
build:
	docker-compose build

# Start the app and database
up:
	docker-compose up

# Stop all running containers
down:
	docker-compose down

# Restart everything
restart: down up

# Show logs
logs:
	docker-compose logs -f
run:
	go run main.go