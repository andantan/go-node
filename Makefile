build:
	@go build -o ./bin/node/node

run: build
	./bin/node/node

test-debug:
	@go test -v ./...

test:
	@go test ./...