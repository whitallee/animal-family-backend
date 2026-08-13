# ENVIRONMENT has no default (config/env.go fails fast if it's missing or
# unrecognized) so every target that execs Go code sets it explicitly here.
# ENVIRONMENT=development is the only value that makes the binary load .env
# itself via godotenv - it handles quoted values correctly. Targets that need
# .env at the shell/Make level instead (test, seed) source it explicitly via
# `.`, since Make's own `include` does NOT strip quotes (KEY="value" would
# literally become "value" with the quotes as part of the string) and go
# test's per-package working directory means godotenv can't find a
# root-level .env on its own.

build:
	@go build -o bin/animal-family-backend cmd/main.go

test:
	@set -a; . ./.env; set +a; ENVIRONMENT=test go test -v ./...

run: build
	@ENVIRONMENT=development ./bin/animal-family-backend

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

seed:
	@set -a; . ./.env; set +a; PGPASSWORD=$$DB_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$DB_USER -d $$DB_NAME -f cmd/migrate/seed/seed.sql

migrate-up:
	@ENVIRONMENT=development go run cmd/migrate/main.go up

migrate-down:
	@ENVIRONMENT=development go run cmd/migrate/main.go down

vapid-keys:
	@go run cmd/vapidgen/main.go

install-hooks:
	@git config core.hooksPath .githooks
	@chmod +x .githooks/*
	@echo "Git hooks installed (gofmt + go vet will run before each commit)"
