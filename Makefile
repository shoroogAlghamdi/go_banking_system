
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server

postgres:
	docker run --name postgres-12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=shoroog -d postgres:12-alpine
createdb:
	docker exec -it postgres-12 createdb --username=root --owner=root banking_system
dropdb:
	docker exec -it postgres-12 dropdb --username=root --owner=root banking_system
migrateup:
	migrate -path db/migration -database "postgresql://root:shoroog@localhost:5432/banking_system?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:shoroog@localhost:5432/banking_system?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
# it is better to name this file as Makefile with small if, when using mac