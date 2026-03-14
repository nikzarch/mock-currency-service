APP_NAME=mock-currency-service

.PHONY: run build test fmt vet tidy up down logs db-reset

run:
	go run ./cmd/mock_service

build:
	go build -o ./bin/$(APP_NAME) ./cmd/mock_service

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

vet:
	go vet ./...

tidy:
	go mod tidy

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f

db-reset:
	docker compose down -v
	docker compose up -d