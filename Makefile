run:
	@go run ./cmd/web


run/migrate:
	@go run ./cmd/web -migrate=true

DB_URL=postgresql://root:secret@localhost:5435/news_app?sslmode=disable

postgres:
	sudo docker run --name postgresNewsApp -p 5435:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	sudo docker exec -it postgresNewsApp createdb --username=root --owner=root news_app

dropdb:
	sudo docker exec -it postgresNewsApp dropdb news_app


.PHONY: run postgres createdb dropdb