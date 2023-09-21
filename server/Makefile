include .env

run:
	@echo "Running server..."
	@go run main.go

build:
	@echo "Building server..."
	@go build -o bin/server main.go

migrate_up:
	@echo "Migrating up..."
	sqlx migrate run --database-url "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"

packages_install:
	@echo "Installing packages..."
	@go mod tidy

packages_update:
	@echo "Updating packages..."
	@go get -u
