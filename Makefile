.PHONY: fmt lint tidy

fmt:
	go fmt ./...
	goimports -w .

lint:
	golangci-lint run ./...

tidy:
	go mod tidy
