DOCKER_COMPOSE := docker-compose
APP_NAME := deathfire-arsenal

.PHONY: build
build:
	@$(DOCKER_COMPOSE) build

.PHONY: run
run:
	@$(DOCKER_COMPOSE) up

.PHONY: clean
clean:
	@$(DOCKER_COMPOSE) down -v

.PHONY: help
help:
	@echo "  build         Build the Go application using Docker Compose."
	@echo "  run           Run the Go application along with MongoDB and Redis using Docker Compose."
	@echo "  clean         Stop and remove the containers, networks, volumes, and images created by Docker Compose."
