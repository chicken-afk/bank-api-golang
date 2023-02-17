stoppostgres:
	docker stop postgres12 

startpostgres:
	docker start postgres12

postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root banks

dropdb:
	docker exec -it postgres12 dropdb --username=root --owner=root banks

createmigration:
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banks?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/banks?sslmode=disable" -verbose down

sqlc:
	docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate
	
test:
	go test -v -cover ./...

.PHONY: stoppostgres startpostgres postgres createdb dropdb migrateup migratedown sqlc test createmigration