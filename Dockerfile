# Stage 1: compile both binaries
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency manifests first so this layer is cached until deps change
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/migrate cmd/migrate/main.go

# Stage 2: lean runtime image (~20MB vs ~800MB for the builder)
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/server      ./bin/server
COPY --from=builder /app/bin/migrate     ./bin/migrate
# Migrations SQL files are read at runtime by the migrate binary
COPY --from=builder /app/cmd/migrate/migrations ./cmd/migrate/migrations

COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
