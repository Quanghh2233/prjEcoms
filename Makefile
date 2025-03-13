.PHONY: up down build migrate schema seed seed-go-pg start test init-db

# Docker commands
up:
	docker-compose up

down:
	docker-compose down

build:
	docker-compose build

# Database operations
init-db:
	go run pkg/db/dbinit/db_init.go

schema:
	go run pkg/db/init/init.go

# Legacy sqlc migrations
migrate:
	migrate -path pkg/db/migration -database "postgresql://qhh:2203@localhost:5432/ecommerce?sslmode=disable" -verbose up

migrate-down:
	migrate -path pkg/db/migration -database "postgresql://qhh:2203@localhost:5432/ecommerce?sslmode=disable" -verbose down

# Seed data
seed-go-pg:
	go run pkg/db/seed/seed_go_pg.go

seed:
	go run pkg/db/seed/seed.go

# Run SQLC code generation
sqlc:
	sqlc generate

# Start application locally
start:
	go run cmd/api/main.go

# Run tests
test:
	go test ./...
