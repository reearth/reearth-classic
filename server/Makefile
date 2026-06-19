lint:
	golangci-lint run --fix

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

run-app-op:
	@command -v op >/dev/null 2>&1 || { echo "Error: 1Password CLI not installed. See https://developer.1password.com/docs/cli/get-started/"; exit 1; }
	@test -f .env.op || { echo "Error: .env.op not found. Create server/.env.op with 1Password references (see PR description)."; exit 1; }
	op run --env-file=.env.op -- go run ./cmd/reearth

gql:
	go generate ./internal/adapter/gql

.PHONY: lint test e2e build run-app run-app-op run-db gql
