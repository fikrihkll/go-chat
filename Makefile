migrate:
	go run ./cmd/cli/... migrate up

migrate-down:
	go run ./cmd/cli/... migrate down

migrate-step-down:
	go run ./cmd/cli/... migrate step-down

migrate-info:
	go run ./cmd/cli/... migrate info

migrate-force:
	go run ./cmd/cli/... migrate force $(or $(version), 1)

seed:
	go run ./cmd/cli/... seed

test:
	go test ./application/chat/tests/... -v

swagger:
	swag init -g cmd/api/main.go

docker:
	docker-compose up -d

docker-down:
	docker-compose down

run:
	go run ./cmd/api/main.go

migration:
	@echo "Running migration"
	chmod +x generate_migration.sh && ./generate_migration.sh

