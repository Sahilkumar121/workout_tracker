.Phony: build run

build:
	@mkdir -p ./bin
	@go build -o ./bin/main ./cmd/api

run: build
	@./bin/main

migrate-up:
	@go run ./cmd/migrate up

migrate-down:
	@go run ./cmd/migrate down