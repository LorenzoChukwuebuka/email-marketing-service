ROOT_DIR := ./backend/
FE_DIR := ./frontend/
DIST_DIR := $(FE_DIR)dist
CLIENT_DIR := $(ROOT_DIR)client

DB_DRIVER := postgres
DB_URL := postgres://postgres:Lorenzo@localhost:5432/email_marketing_service?sslmode=disable
DB_MIGRATIONS_DIR := $(ROOT_DIR)internal/db/migrations

.PHONY: build run buildfe runfe tidy npmi rundev buildfe git docker vet docker-exec docker-restart docker-dev

tidy: 
	cd $(ROOT_DIR) && go mod tidy

npmi:
	cd $(FE_DIR) && npm i $(ARGS)

rundev:
	cd $(FE_DIR) && npm run dev

buildfe:
	cd $(FE_DIR) && npm run build
ifeq ($(OS),Windows_NT)
	powershell -Command "if (Test-Path $(subst /,\,$(CLIENT_DIR))) { Remove-Item -Recurse -Force $(subst /,\,$(CLIENT_DIR)) }"
	powershell -Command "New-Item -ItemType Directory -Force -Path $(subst /,\,$(CLIENT_DIR))"
	powershell -Command "Copy-Item -Path '$(subst /,\,$(DIST_DIR))\*' -Destination '$(subst /,\,$(CLIENT_DIR))' -Recurse -Force"
else
	rm -rf $(CLIENT_DIR)
	mkdir -p $(CLIENT_DIR)
	cp -r $(DIST_DIR)/* $(CLIENT_DIR)/
endif

deploy:
ifndef MSG
	$(error MSG is required. Usage: make deploy MSG="your commit message")
endif
	git add .
	git commit -m "$(MSG)"
	git push
	
docker:
	docker-compose up --build -d

vet:
	cd $(ROOT_DIR) && go vet

docker-exec:
	docker exec -it $(ARGS) bash

docker-restart:
	docker-compose stop
	docker-compose down
	$(MAKE) docker

docker-dev:
	docker-compose -f compose.dev.yaml up --build -d

docker-staging:
	docker-compose -f compose.staging.yaml up --build -d

run:
	cd $(ROOT_DIR) && go run cmd/api/main.go

sqlc-gen:
	cd $(ROOT_DIR) && sqlc generate

# ==================================================================================== #
# SQL MIGRATIONS
# ==================================================================================== #

.PHONY: create-migration migrate-force migrate-up migrate-down migrate-rollback migrate-drop migrate-to

create-migration: ## Create a new migration file for a table, e.g., make create-migration table_name=bug_report
ifdef table_name
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -ext sql -dir $(DB_MIGRATIONS_DIR) -seq $(table_name)
else
	@echo "Please provide a table_name argument, e.g., make create-migration table_name=bug_report"
endif

migrate-force: ## Force migration to specific version
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR) force $(version)

migrate-up: ## Migrate the database schema up to the latest version
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR) up

migrate-down: ## Rollback the database schema by one migration
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR) down

migrate-rollback: ## Rollback the database schema by specified steps
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR) down $(step)

migrate-drop: ## Drop all migrations (WARNING: This will delete all data)
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR) drop -f

migrate-to: ## Migrate the database schema to a specific version, e.g., make migrate-to version=1
ifdef version
	go run -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest -database=$(DB_URL) -path=$(DB_MIGRATIONS_DIR) goto $(version)
else
	@echo "Please provide a version argument, e.g., make migrate-to version=1"
endif


debug-paths:
	@echo "ROOT_DIR: $(ROOT_DIR)"
	@echo "DB_MIGRATIONS_DIR: $(DB_MIGRATIONS_DIR)"
	@echo "Current directory: $(shell pwd)"