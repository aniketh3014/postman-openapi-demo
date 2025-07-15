build:
	@go build -o bin/server ./cmd/server/

run: build
	@./bin/server
