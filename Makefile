postgres_container:
	docker run --name postgres16 --network=bank-network -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

create_db:
	docker exec -it postgres16 createdb --username=root --owner=root simplebank

drop_db:
	docker exec -it postgres16 dropdb simplebank

# migrate_init:
# 	migrate create -ext sql -dir db/migrations -seq init_schema

up_migrate:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simplebank?sslmode=disable" -verbose up

down_migrate:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simplebank?sslmode=disable" -verbose down

up_migrate_last:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simplebank?sslmode=disable" -verbose up 1

down_migrate_last:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simplebank?sslmode=disable" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)
sqlc:
	sqlc generate


test:
	go test -v -cover ./...

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/komron-dev/bank/db/sqlc Store

gen_proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
        --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank \
        proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl cli --package pb --service Bank
run:
	go run main.go
.PHONY: postgres_container create_db drop_db up_migrate up_migrate_last down_migrate down_migrate_last test run mock gen_proto