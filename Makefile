# build/run/migrate-up/migrate-down/vapid-keys all exec a Go binary from the
# repo root, which loads .env itself via godotenv (config/env.go) and handles
# quoted values correctly. Only targets below that need .env at the shell/Make
# level (test, seed) source it explicitly via `.`, since Make's own `include`
# does NOT strip quotes (KEY="value" would literally become "value" with the
# quotes as part of the string) and go test's per-package working directory
# means godotenv can't find a root-level .env on its own.

build:
	@go build -o bin/animal-family-backend cmd/main.go

test:
	@set -a; . ./.env; set +a; ENVIRONMENT=test go test -v ./...

run: build
	@./bin/animal-family-backend

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

seed:
	@set -a; . ./.env; set +a; PGPASSWORD=$$DB_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$DB_USER -d $$DB_NAME -f cmd/migrate/seed/seed.sql

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

vapid-keys:
	@go run cmd/vapidgen/main.go
	