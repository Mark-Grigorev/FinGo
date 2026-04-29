BIN     := fingo
CMD     := ./cmd
SWAG    := $(shell go env GOPATH)/bin/swag

.PHONY: run build test tidy swagger up down logs ps

run:
	go run $(CMD)

build:
	go build -o $(BIN) $(CMD)

test:
	go test ./...

tidy:
	go mod tidy

swagger:
	$(SWAG) init -g $(CMD)/main.go -o ./docs

up:
	docker compose up -d --build

down:
	docker compose down

logs:
	docker compose logs -f app

ps:
	docker compose ps
