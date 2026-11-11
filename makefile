.Phony: build run

build:
	@mkdir -p ./bin
	@go build -o ./bin/main ./cmd/api

run: build
	@./bin/main