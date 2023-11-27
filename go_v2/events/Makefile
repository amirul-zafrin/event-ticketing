DB_URL=postgresql://root:secret@localhost:5432/events?sslmode=disable

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it postgres psql -U root events

dropdb:
	docker exec -it postgres dropdb events

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

sqlc:
	sqlc generate
phony: postgres createdb dropdb migrateup sqlc