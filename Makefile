include .env
export

.PHONY: run build test test-coverage docker-up docker-down docker-logs migrate-create migrate-up migrate-down seed act-build act-lint act-test act-clean act-build-local


run:
	go run ./cmd/server/main.go

build:
	go build -o bin/server ./cmd/server

test:
	go test -v -race ./...

test-coverage:
	go test -coverprofile=coverage.out ./... || true
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-up:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-down:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down

seed:
	go run ./seeds/seed.go

act-build-local:
	act -j build -W .github/workflows/docker-local.yml --secret-file .env.act

act-lint:
	act -j lint --secret-file .env.act

act-test:
	act -j test --secret-file .env.act