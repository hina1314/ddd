code:
	goctl api go --api ./dsl/user.api --dir .

mysql:
	goctl model mysql datasource --url="root:root@tcp(192.168.1.10:5432)/yajiang" -table="*" -dir="./internal/model"

postgres:
	goctl model pg datasource --url="postgres://postgres:mysecretpassword@192.168.1.10:5432/postgres" --table="*" -dir="./db/model"

migrate_init:
	migrate create -ext sql -dir ./db/migration -seq init_schema

migrate_up:
	migrate -path db/migration -database "postgres://postgres:mysecretpassword@192.168.1.10:5432/postgres?sslmode=disable" --verbose up

migrate_down:
	migrate -path db/migration -database "postgres://postgres:mysecretpassword@192.168.1.10:5432/postgres?sslmode=disable" --verbose down

sqlc:
	sqlc generate
wire:
	wire gen ./internal/di/wire.go
.PHONY: code migrate_up migrate_down sqlc