postgres = "postgres://postgres:123456@127.0.0.1:5432/postgres?sslmode=disable"

migrate_init:
	migrate create -ext sql -dir ./db/migration -seq $(name)

migrate_up:
	migrate -path db/migration -database $(postgres) --verbose up

migrate_down:
	migrate -path db/migration -database $(postgres) --verbose down

sqlc:
	sqlc generate
wire:
	wire gen ./internal/di/wire.go
.PHONY: code migrate_up migrate_down sqlc