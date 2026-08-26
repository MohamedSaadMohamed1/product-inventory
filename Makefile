.PHONY: run test test-race docker-up docker-down sqlc lint db-backup db-restore

run:
	go run cmd/api/main.go

test:
	go test -v ./...

test-race:
	go test -v -race ./...

docker-up:
	docker compose up -d

docker-down:
	docker compose down

sqlc:
	sqlc generate

lint:
	golangci-lint run

db-backup:
	docker exec -t inventory_postgres pg_dump -U postgres -d inventory_db > database/backup.sql

db-restore:
	docker exec -i inventory_postgres psql -U postgres -d inventory_db < database/backup.sql
