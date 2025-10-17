postgres:
	docker run -d --name postgres-bpay -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -v ~/docker/postgres-amole:/var/lib/postgresql/data -p 5432:5432 postgres:17-alpine
createdb:
	docker exec -it postgres-bpay createdb --username=root --owner=root master-db
dropdb:
	docker exec -it postgres-bpay dropdb master-db
migrateup:
	migrate -path script/migrations -database "postgresql://admin:admin@localhost:5432/master-db?sslmode=disable" -verbose up
migratedown:
	migrate -path script/migrations -database "postgresql://admin:admin@localhost:5432/master-db??sslmode=disable" -verbose down

.PHONY: createdb dropdb migrateup migratedown