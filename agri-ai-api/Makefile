.PHONY: run db-up db-down migrate-up migrate-down swagger

run: swagger
	go run cmd/server/main.go

db-up:
	docker-compose up -d

db-down:
	docker-compose down

migrate-up:
	migrate -path migrations -database "postgres://admin:password@localhost:5432/agri_ai?sslmode=disable" -verbose up

migrate-down:
	migrate -path migrations -database "postgres://admin:password@localhost:5432/agri_ai?sslmode=disable" -verbose down

swagger:
	swag init -g cmd/server/main.go
