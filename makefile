DB_SOURCE ?=
SQLC_VERSION := v1.31.1
WIRE_VERSION := v0.6.0

.PHONY: generate sqlc wire test vet fmt check migrate_init migrate_up migrate_down

generate: sqlc wire

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) generate

wire:
	go run github.com/google/wire/cmd/wire@$(WIRE_VERSION) ./internal/di

fmt:
	gofmt -w .

test:
	go test ./...

vet:
	go vet ./...

check: generate fmt vet test
	git diff --exit-code -- db/model internal/di/wire_gen.go

migrate_init:
	@test -n "$(name)" || (echo "name is required" && exit 1)
	migrate create -ext sql -dir ./db/migration -seq "$(name)"

migrate_up:
	@test -n "$(DB_SOURCE)" || (echo "DB_SOURCE is required" && exit 1)
	migrate -path db/migration -database "$(DB_SOURCE)" --verbose up

migrate_down:
	@test -n "$(DB_SOURCE)" || (echo "DB_SOURCE is required" && exit 1)
	migrate -path db/migration -database "$(DB_SOURCE)" --verbose down 1
