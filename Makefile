migrateup:
	migrate -path postgres/migrations -database "postgres://postgres:Veg@zr01@103.127.99.34:5432/gomaster?sslmode=disable" -verbose up

migratedown:
	migrate -path postgres/migrations -database "postgresCont://postgres:Veg@zr01@103.127.99.34:5432/gomaster?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

migration:
	@read -p "enter migration name : " name; \
		migrate create -ext sql -dir postgres/migrations $$name

.PHONY: migrateup, migratedown, sqlc, test, migration
