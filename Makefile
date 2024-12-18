build:
	@go build -o bin/animal-family-backend cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/animal-family-backend