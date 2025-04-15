tidy:
	go mod tidy

build: 
	mkdir -p bin
	go build -o bin/federator main.go

start: build
	./bin/federator

dev-up:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.dev.yml \
		--env-file .env \
		up -d

dev-down:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.dev.yml \
		--env-file .env \
		down -v

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