
BINARY_NAME := pr-reviewer

build:
	@go build -o bin/$(BINARY_NAME) cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/app

debug:
	@dlv debug --headless --listen=:2345 --log --api-version=2 ./cmd/main.go

# Clean up generated files
clean:
	rm -f bin/$(BINARY_NAME)
