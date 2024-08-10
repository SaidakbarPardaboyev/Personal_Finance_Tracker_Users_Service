-include .env
export

CURRENT_DIR=$(shell pwd)
DATABASE_URL = postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DBNAME}?sslmode=disable

# run service
.PHONY: run
run:
	go run cmd/main.go

# migrate
.PHONY: migrate
migrate:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DBNAME}?sslmode=disable up

DB_URL := "postgres://postgres:123@localhost:5432/db?sslmode=disable"

tidy:
	@go mod tidy

mig-create:
	@if [ -z "$(name)" ]; then \
	  read -p "Enter migration name: " name; \
	fi; \
	migrate create -ext sql -dir migrations -seq $$name

mig-up:
	@migrate -database "$(DATABASE_URL)" -path migrations up

mig-down:
	@migrate -database "$(DATABASE_URL)" -path migrations down

mig-force:
	@if [ -z "$(version)" ]; then \
	  read -p "Enter migration version: " version; \
	fi; \
	migrate -database "$(DATABASE_URL)" -path migrations force $$version


pull-proto-module:
	git submodule update --init --recursive

update-proto-module:
	git submodule update --remote --merge

swag-gen:
	swag init -g ./api/router.go -o api/docs

proto-gen:
	@./scripts/gen-proto.sh $(CURRENT_DIR)

test:
	@go test ./storage/postgres
