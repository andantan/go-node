build:
	@go build -o ./bin/node/node

run: build
	./bin/node/node

test:
	@go test -v ./...