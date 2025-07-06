# Variables
DOCKER_IMAGE_NAME = mankings/federator
DOCKER_TAG ?= 1.5
DOCKERFILE_PATH = deployment/docker/Dockerfile

tidy:
	go mod tidy

build: 
	mkdir -p bin
	go build -o bin/federator main.go

start: build
	./bin/federator

up:
	./bin/federator

# Docker build rules
docker-build:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) -f $(DOCKERFILE_PATH) .

docker-build-latest: docker-build

docker-tag:
	docker tag $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) $(DOCKER_IMAGE_NAME):latest

docker-push:
	docker push $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

docker-push-latest:
	docker push $(DOCKER_IMAGE_NAME):latest

docker-build-push: docker-build docker-push

docker-build-push-latest: docker-build docker-tag docker-push-latest

half-up:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.half.yml \
		--env-file .env.half \
		up -d \

half-down:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.half.yml \
		--env-file .env.half \
		down -v

half-build:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.half.yml \
		--env-file .env.half \
		up -d --build

full-up:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.full.yml \
		--env-file .env.full \
		up -d

full-down:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.full.yml \
		--env-file .env.full \
		down -v

full-build:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.full.yml \
		--env-file .env.full \
		up -d --build

clean:
	rm -rf bin/*
	docker rmi $(APP_NAME) -f
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.full.yml \
		down -v
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.half.yml \
		down -v
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.dev.yml \
		down -v