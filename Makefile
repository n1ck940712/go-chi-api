DOCKER_NETWORK = default_network
DOCKER_YML = build/docker/docker-compose.yml
DOCKER_COMPOSE_COMMAND = docker-compose -f $(DOCKER_YML) --project-directory ./

up:
	$(DOCKER_COMPOSE_COMMAND) up -d

down:
	$(DOCKER_COMPOSE_COMMAND) down
	rm -f main

rebuild:
	$(DOCKER_COMPOSE_COMMAND) up -d --build

restart:
	$(DOCKER_COMPOSE_COMMAND) restart

network_create:
	docker network create $(DOCKER_NETWORK)

network_remove:
	docker network rm $(DOCKER_NETWORK)
