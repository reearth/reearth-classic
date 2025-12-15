lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run --timeout=60m

test:
	go test -race -short -v ./...

e2e:
	go test -v ./e2e/...

build:
	go build ./cmd/reearth

run-app:
	go run ./cmd/reearth

run-db:
	docker compose -f ../docker-compose.yml up -d reearth-mongo

gql:
	go generate ./internal/adapter/gql

.PHONY: lint test e2e build run-app run-db gql
