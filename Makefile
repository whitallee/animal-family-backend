-include .env
export

build:
	@go build -o bin/animal-family-backend cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/animal-family-backend

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

seed:
	@PGPASSWORD=${DB_PASSWORD} psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} -f cmd/migrate/seed/seed.sql

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

vapid-keys:
	@go run cmd/vapidgen/main.go
	