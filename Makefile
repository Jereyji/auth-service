DOCKER_COMPOSE_PATH=$(CURDIR)/deployments/docker-compose.yaml

create-env:
	bash scripts/create_env.sh

up-service:
	docker-compose -f $(DOCKER_COMPOSE_PATH) up --build

down-service:
	docker-compose -f $(DOCKER_COMPOSE_PATH) down
# docker-compose -f $(DOCKER_COMPOSE_PATH) down --volumes
