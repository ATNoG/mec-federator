include .env

tidy:
	go mod tidy

build: 
	mkdir -p bin
	go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

start: build
	./bin/$(APP_NAME)

docker-build:
	docker build -t $(APP_NAME) . -f deployment/docker/Dockerfile

docker-start: docker-build
	docker run -p $(API_PORT):$(API_PORT) $(APP_NAME)

dc-up:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.op-a.yml \
		--env-file .env.a \
		up -d \
		$(if $(BUILD),--build)

dc-down:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.op-a.yml \
		down -v

db-up:
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.op-a.yml \
		--env-file .env.a \
		up -d \
		mongo mongo-express

clean:
	rm -rf bin/*
	docker rmi $(APP_NAME) -f
	docker compose \
		--project-directory . \
		-f deployment/docker/docker-compose.op-a.yml \
		down -v
