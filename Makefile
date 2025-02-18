include .env

tidy:
	go mod tidy

build: 
	go build -o bin/$(APP_NAME) cmd/$(APP_NAME)/main.go

start:
	./bin/$(APP_NAME)

restart: build start