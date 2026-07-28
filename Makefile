ifneq (,$(wildcard .env))
include .env
export
endif

.PHONY: fmt lint tidy migrate-up migrate-down

fmt:
	go fmt ./...
	goimports -w .

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

migrate-up:
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir db/migrations postgres "$(DATABASE_URL)" down
