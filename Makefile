postgres_container:
	docker run --name postgres16 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

create_db:
	docker exec -it postgres16 createdb --username=root --owner=root simple_bank

drop_db:
	docker exec -it postgres16 dropdb simple_bank

# migrate_init:
# 	migrate create -ext sql -dir db/migrations -seq init_schema

up_migrate:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up

down_migrate:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

# creates 3 accounts on the database for testing
pre_test:
	go test -v -cover -run 'CreateAccount' ./...
	go test -v -cover -run 'CreateAccount' ./...
	go test -v -cover -run 'CreateAccount' ./...

test:
	go test -v -cover ./...

run:
	go run main.go
.PHONY: postgres_container create_db drop_db up_migrate down_migrate test run