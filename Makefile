MIGRATIONS_DIR=./migrations
ENV_FILE=.env

include $(ENV_FILE)
export

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: migrate-create migrate-up migrate-down server-start test-api lint-api

server-start:
	@go run ./cmd

migrate-create:
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down -1

test-api:
	@go test ./... -v -cover

lint-api:
	golangci-lint run