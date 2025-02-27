include .env

tidy:
	go mod tidy

build: 
	go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

start:
	./bin/$(APP_NAME)


docker:
	docker build -t $(APP_NAME) . -f deployment/docker/Dockerfile