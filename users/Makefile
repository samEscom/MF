SHELL=/bin/bash
DB_USER := $(shell grep DB_USER .env.local | cut -d '=' -f 2)
DB_PASSWORD := $(shell grep DB_PASSWORD .env.local | cut -d '=' -f 2)
DB_HOST := $(shell grep DB_HOST .env.local | cut -d '=' -f 2)
DB_PORT := $(shell grep DB_PORT .env.local | cut -d '=' -f 2)
DB_NAME := $(shell grep DB_NAME .env.local | cut -d '=' -f 2)

.PHONY: download
download:
	go mod download
	go mod tidy

.PHONY: run-local
run-local:
	docker-compose -f infra/local/docker-compose.yml up -d


.PHONY: migrateup
migrateup:
	migrate -path ./db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose up

.PHONY: migratedown
migratedown:
	migrate -path ./db/migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" -verbose down

.PHONY: watch-coverage
watch-coverage:
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out


.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint -v run

.PHONY: mocks
mocks:
	rm -fr ./mocks
	mockery --config mockeryconfig.yml
