include .env

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  run              	Run server"
	@echo "  build            	Build Binaries for linux, windows and mac"
	@echo "  run_build        	Run server from build"
	@echo "  test             	Run tests"
	@echo "  lint             	Run linter"
	@echo "  migrate_new      	Create new migration"
	@echo "  migrate_up       	Apply pending migrations"
	@echo "  packages_install 	Install packages"
	@echo "  packages_update  	Update packages"
	@echo "  dev              	Run dev server"
	@echo "  swagger_docgen   	Generate Swagger Docs"
	@echo "  docgen           	Generate OpenAPIv3 Docs and Swagger Docs"
	@echo "  cloc             	Count lines of code"

run:
	@echo "Running server..."
	@go run main.go

build:
	@echo "Building server..."
	@GOOS=linux GOARCH=amd64 go build -o bin/server main.go 
	@GOOS=windows GOARCH=amd64 go build -o bin/server.exe main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/server.darwin main.go
	@echo "Done building for linux, windows and mac (x64 only)"
	@echo "Binaries are in bin/ directory"

run_build:
	@echo "Running server..."
	@./bin/server

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run

migrate_new:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	sqlx migrate add $$name 

migrate_up:
	@echo "Migrating up..."
	sqlx migrate run --database-url "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}"

packages_install:
	@echo "Installing packages..."
	@go mod tidy

packages_update:
	@echo "Updating packages..."
	@go get -u

dev:
	@air

swagger_docgen:
	@echo "Generating Swagger Docs..."
	@echo ""
	@rm -rf docs/
	@swag init
	@echo ""
	@echo "Done generating docs."
	@echo ""
	@echo "Visit: http://localhost:${PORT}/swagger/index.html"

docgen: swagger_docgen
	@echo ""
	@echo "Generating OpenAPIv3 Docs..."
	@echo ""
	@rm -rf docs/openapi.yaml
	@npx -p swagger2openapi swagger2openapi --yaml --outfile docs/openapi.yaml "http://localhost:${PORT}/swagger/doc.json"
	@echo "Done generating OpenAPIv3 Docs."
