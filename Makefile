include .env
export

.PHONY: run build test test-coverage docker-up docker-down docker-logs migrate-create migrate-up migrate-down

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

# Start all containers (PostgreSQL, migrations, seed, app) in order:
# 1. postgres - starts and passes healthcheck
# 2. migrate - runs migrations, exits
# 3. seed - populates database with test data, exits
# 4. app - starts after seed completes successfully
docker-up:
	docker-compose up -d

# Stop all containers
docker-down:
	docker-compose down

# View logs from all containers
docker-logs:
	docker-compose logs -f

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

migrate-up:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" up

migrate-down:
	migrate -path migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}" down