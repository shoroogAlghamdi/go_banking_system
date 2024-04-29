
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migratedown1 migrateup1 dockercontainer

postgres:
	docker run --name postgres-12 --network=bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=shoroog -d postgres:12-alpine
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
mock: 
	mockgen -package mockdb -destination db/mock/store.go github.com/shoroogAlghamdi/banking_system/db/sqlc Store	
migratedown1:
	migrate -path db/migration -database "postgresql://root:shoroog@localhost:5432/banking_system?sslmode=disable" -verbose down 1

migrateup1:
	migrate -path db/migration -database "postgresql://root:shoroog@localhost:5432/banking_system?sslmode=disable" -verbose up 1
dockerimage:
	docker build -t banking_system:latest .          
dockercontainer:
	docker run --name banking_system --network=bank-network -p 8080:8080 -e GIN_MODE=release -e DB_SOURCE=postgresql://root:shoroog@postgres-12:5432/banking_system?sslmode=disable banking_system:latest  
inspectcontainer:
	docker container inspect postgres-12
	docker network inspect fa8c322eb9ea 
userdefineddockernetwork:
	docker network create bank-network
connectcontainertonetwork:
	docker network connect bank-network postgres-12
# it is better to name this file as Makefile with small if, when using mac