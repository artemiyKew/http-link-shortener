.PHONY: build

build: 
	go build -v ./cmd/app

.PHONY: test

run: build
	./app

test: 
	go test -v -race -timeout 30s ./...

compose:
	docker compose up

.DEFAULT_GOAL = build