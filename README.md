# Animal Family Backend

A Go backend API for managing animal care tasks, enclosures, and related data.

**Frontend Repository:** [animal-family-web](https://github.com/whitallee/animal-family-web)

> **Note:** This is currently in Pre-Alpha. If you'd like to collaborate, please reach out! Find my contact info on [whitcodes.dev/contact](https://whitcodes.dev/contact).

## Prerequisites

- Go 1.24 or later
- PostgreSQL
- Make (optional, for using Makefile commands)

## Running Locally

### 1. Clone the Repository

```bash
git clone https://github.com/whitallee/animal-family-backend.git
cd animal-family-backend
```

### 2. Set Up Environment Variables

Copy `.env.example` to `.env` and fill in real values:

```bash
cp .env.example .env
```

See `.env.example` for the full list of variables and their defaults.

### 3. Set Up Database

Ensure PostgreSQL is running and create the database:

```bash
createdb animal_family
```

### 4. Run Migrations

Run database migrations and seed data:

```bash
make migrate-up
```

### 5. Build and Run

Using Make:

```bash
make run
```

Or manually:

```bash
go build -o bin/animal-family-backend cmd/main.go
./bin/animal-family-backend
```

The API server will start on `http://localhost:8080`.

### 6. Install Git Hooks (one-time)

```bash
make install-hooks
```

Runs `gofmt` + `go vet` before each commit.

## Available Make Commands

- `make build` - Build the application
- `make run` - Build and run the application
- `make test` - Run all tests
- `make migrate-up` - Run database migrations and seed data
- `make migrate-down` - Rollback database migrations
- `make seed` - Seed the database with initial data
- `make migration <name>` - Create a new migration file
- `make install-hooks` - Install git hooks (gofmt + go vet before each commit)
- `make spec` - Regenerate the v2 OpenAPI contract (`docs/openapi.json`)
- `make spec-check` - Verify the committed contract matches the code

## API Versions

`/api/v1` is the original API. `/api/v2` is a smaller, REST-conventional
replacement being built alongside it — plural collections, IDs in the path, and
no separate `/admin/*` route tree. Both are served; v1 will be retired once the
frontend has fully migrated.

The v2 contract is generated from annotations on the handlers into
[`docs/openapi.json`](docs/openapi.json), which is also served at
`GET /openapi.json`. The frontend generates its API client from it.

## Project Structure

- `cmd/` - Application entry points (main.go, api/, migrate/)
- `config/` - Configuration management
- `db/` - Database connection and setup
- `service/` - Business logic and route handlers
- `types/` - Type definitions
- `utils/` - Utility functions

## Database Schema

[Entity Relationship Diagram](https://docs.google.com/drawings/d/1Vi1yngr4CeXXt-slRGJsLI35_R-y-oIHlZ466be_wx8/edit?usp=sharing)

## Development Status

- ✅ User
- ✅ Habitat
- ✅ Species
- ✅ Enclosure
- ✅ Animal
- 🔄 Task (in progress)

See [TODO.md](TODO.md) for planned features and known technical debt.
