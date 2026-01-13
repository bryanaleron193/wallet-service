include .env

DATABASE_URL := $(subst \r,,$(DATABASE_URL))
MIGRATIONS_PATH := ./migrations

.PHONY: migration migrate-up migrate-down migrate-force seed

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

%:
	@: