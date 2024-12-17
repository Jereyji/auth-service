DOCKER_COMPOSE_PATH=$(CURDIR)/deployments/docker-compose.yaml

create-env:
	bash scripts/create_env.sh

build-service:
	docker-compose -f $(DOCKER_COMPOSE_PATH) build 

up-service:
	docker-compose -f $(DOCKER_COMPOSE_PATH) up
