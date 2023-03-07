.DEFAULT_GOAL := build

test: lint unit integration ## Run all tests

prepare: ## Prepare project: install dependencies
	go install github.com/jackc/tern@latest

lint: ## Run the linter
	golangci-lint run --build-tags=integration

unit: ## Run unit tests
	go test -race -count=1 ./...

integration: reset-db fixtures ## Run integration tests
	go test -race -p=1 -count=1 --tags=integration ./...

migrate-up: ## Migrate DB schema to newer version
	tern migrate -c database/migrations/tern.conf -m database/migrations

reset-schema: ## Recreates the "public" schema
	PGPASSWORD=todo psql -v ON_ERROR_STOP=1 -U todo -h localhost todo -c "DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;"

reset-db: reset-schema migrate-up ## Reset database schema

run-server:## Run the server on the host
	go run cmd/server/main.go

build: test
	go build -o dist/local/server cmd/server/main.go

fixtures:
	PGPASSWORD=todo psql -v ON_ERROR_STOP=1 -U todo -h localhost todo -f database/fixtures.sql
