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

up_migrate_last:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up 1

down_migrate_last:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)
sqlc:
	sqlc generate


test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/komron-dev/bank/db/sqlc Store

run:
	go run main.go
.PHONY: postgres_container create_db drop_db up_migrate up_migrate_last down_migrate down_migrate_last test run mock