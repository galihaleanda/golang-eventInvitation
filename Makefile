.PHONY: run build docker-up docker-down migrate tidy

run:
	go run ./cmd/api

build:
	go build -o bin/server ./cmd/api

tidy:
	go mod tidy

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	psql "$(DB_URL)" -f migrations/0001_init_schema.up.sql

migrate-down:
	psql "$(DB_URL)" -f migrations/0001_init_schema.down.sql

# Example: make migrate-up DB_URL="postgres://postgres:postgres@localhost:5432/event_invitation?sslmode=disable"
