PROJECT_NAME=auth_service_test
SERVICE_PATH=$(CURDIR)/cmd/main.go
DEPLOYMENTS_PATH=$(CURDIR)/deployments
DOCKER_COMPOSE_PATH=$(DEPLOYMENTS_PATH)/docker-compose.yaml

.PHONY: env
env:
	bash scripts/create_env.sh

.PHONY: build-service
build-service:
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) build

.PHONY: up-service
up-service:
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) up -d

.PHONY: down-service
down-service:
	docker-compose -p $(PROJECT_NAME) -f $(DOCKER_COMPOSE_PATH) down

.PHONY: run-service
run-service:
	go run $(SERVICE_PATH)

.PHONY: run-functional-tests
run-functional-tests:
	go test tests/functional_test.go