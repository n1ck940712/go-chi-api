DOCKER_NETWORK = default_network
DOCKER_YML = build/docker/docker-compose.yml
DOCKER_COMPOSE_COMMAND = docker-compose -f $(DOCKER_YML) --project-directory ./

help:
	@sed -ne '/@sed/!s/## //p'  $(MAKEFILE_LIST)

up: ## Create and start containers
	$(DOCKER_COMPOSE_COMMAND) up -d

down: ## Stop and remove resources
	$(DOCKER_COMPOSE_COMMAND) down
	rm -f main

stop: ## Stop containers
	$(DOCKER_COMPOSE_COMMAND) stop

start: ## Start containers
	$(DOCKER_COMPOSE_COMMAND) start

rebuild: ## Rebuild docker containers
	$(DOCKER_COMPOSE_COMMAND) up -d --build

restart: ## Restart docker containers
	$(DOCKER_COMPOSE_COMMAND) restart

psql: ## Connect to postgres container
	$(DOCKER_COMPOSE_COMMAND) exec postgres_db psql -U postgres

hotreload_on: ## enable hotreload
	touch .env
	echo 'DOCKERFILE_BASE="local.hotreload"' >> .env

hotreload_off: ## disable hotreload
	rm -f .env

network_create: ## Create default docker network
	docker network create $(DOCKER_NETWORK)

network_remove: ## Remove default docker network
	docker network rm $(DOCKER_NETWORK)
