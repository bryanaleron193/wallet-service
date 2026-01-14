include .env

DATABASE_URL := $(shell powershell -Command "$$url = '$(DATABASE_URL)'.Replace('\r', '').Replace('@db:', '@127.0.0.1:').Replace(':5432/', ':$(DB_PORT)/'); Write-Output $$url")
MIGRATIONS_PATH := ./migrations

.PHONY: migration migrate-up migrate-down migrate-force seed gen-docs

migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" up

migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" down $(filter-out $@,$(MAKECMDGOALS))

migrate-force:
	@migrate -path=$(MIGRATIONS_PATH) -database="$(DATABASE_URL)" force $(filter-out $@,$(MAKECMDGOALS))

seed:
	@go run ./cmd/seed/main.go --db "$(DATABASE_URL)"

gen-docs:
	@swag init -g ./api/main.go -d cmd,internal && swag fmt

%:
	@: