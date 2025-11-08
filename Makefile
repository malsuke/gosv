.PHONY: lint test fmt

lint:
	golangci-lint run --config .golangci.yml

test:
	go test ./...

fmt:
	go fmt ./...
